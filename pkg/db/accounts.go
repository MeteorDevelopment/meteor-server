package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/dboslee/lru"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"meteor-server/pkg/core"
	"meteor-server/pkg/discord"
	"net/http"
	"time"
)

type Account struct {
	ID ksuid.KSUID `bson:"id" json:"id"`

	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password []byte `bson:"password" json:"-"`

	Admin   bool `bson:"admin" json:"admin"`
	Donator bool `bson:"donator" json:"donator"`

	DiscordID string `bson:"discord_id" json:"discord_id"`

	MaxMcAccounts int         `bson:"max_mc_accounts" json:"max_mc_accounts"`
	McAccounts    []uuid.UUID `bson:"mc_accounts" json:"mc_accounts"`

	Cape              string `bson:"cape" json:"cape"`
	CanHaveCustomCape bool   `bson:"can_have_custom_cape" json:"can_have_custom_cape"`
}

var usernameCache = lru.NewSync[ksuid.KSUID, string](lru.WithCapacity(20))

var (
	AccountCount int64
	DonatorCount int64
)

func initAccounts() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var err error

	AccountCount, err = accounts.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Failed to get the number of accounts: %s\n", err)
	}

	DonatorCount, err = accounts.CountDocuments(ctx, bson.M{"donator": true})
	if err != nil {
		fmt.Printf("Failed to get the number of donators: %s\n", err)
	}
}

func NewAccount(username string, email string, password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = accounts.InsertOne(nil, Account{
		ID: ksuid.New(),

		Username: username,
		Email:    email,
		Password: pass,

		Admin:   false,
		Donator: false,

		DiscordID: "",

		MaxMcAccounts: 1,
		McAccounts:    []uuid.UUID{},

		Cape:              "",
		CanHaveCustomCape: false,
	})

	return err
}

func GetAccount(r *http.Request) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"id": core.GetAccountID(r)}).Decode(&acc)

	return acc, err
}

func GetAccountId(id ksuid.KSUID) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"id": id}).Decode(&acc)

	return acc, err
}

func GetAccountDiscordId(id string) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"discord_id": id}).Decode(&acc)

	return acc, err
}

func GetAccountUuid(uuid uuid.UUID) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"mc_accounts": uuid}).Decode(&acc)

	return acc, err
}

func GetAccountsWithCape() []Account {
	cursor, err := accounts.Find(nil, bson.M{"cape": bson.M{"$ne": ""}})
	if err != nil {
		log.Fatal(err)
	}

	var acc []Account
	err = cursor.All(nil, &acc)
	if err != nil {
		log.Fatal(err)
	}

	return acc
}

func GetAccountWithUsername(username string) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"username": username}).Decode(&acc)

	return acc, err
}

func GetAccountWithEmail(email string) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"email": email}).Decode(&acc)

	return acc, err
}

func GetAccountWithUsernameOrEmail(name string) (Account, error) {
	acc, err := GetAccountWithUsername(name)

	if err != nil {
		acc, err = GetAccountWithEmail(name)
	}

	return acc, err
}

func GetAccountUsername(id ksuid.KSUID) (string, error) {
	username, ok := usernameCache.Get(id)
	if ok {
		return username, nil
	}

	account, err := GetAccountId(id)
	if err != nil {
		return "", err
	}

	usernameCache.Set(id, account.Username)
	return account.Username, nil
}

func (acc *Account) PasswordMatches(password string) bool {
	return bcrypt.CompareHashAndPassword(acc.Password, []byte(password)) == nil
}

func (acc *Account) LinkDiscord(id string) error {
	// Make sure an account with this discord id doesn't exist
	_, err := GetAccountDiscordId(id)
	if err == nil {
		return errors.New("Discord account already linked.")
	}

	// Put discord ID in database
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"discord_id": id}})

	// If the account has donator but not donator role, give it
	if acc.Donator {
		discord.AddRole(id, discord.DonorRole)
		discord.SendDonorMsg(id)
	}

	// Add account role regardless
	discord.AddRole(id, discord.AccountRole)
	return nil
}

func (acc *Account) UnlinkDiscord() {
	// Remove account related roles
	discord.RemoveRole(acc.DiscordID, discord.AccountRole)
	discord.RemoveRole(acc.DiscordID, discord.DonorRole)

	// Remove discord ID from database
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"discord_id": ""}})
}

func (acc *Account) AddMcAccount(id uuid.UUID) error {
	// Check maximum number of Minecraft accounts
	if len(acc.McAccounts) >= acc.MaxMcAccounts {
		return errors.New("Exceeded maximum number of Minecraft accounts.")
	}

	// Check for duplicate Minecraft accounts
	for _, mcAccount := range acc.McAccounts {
		if mcAccount == id {
			return errors.New("Account already has that Minecraft account linked.")
		}
	}

	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$push": bson.M{"mc_accounts": id.String()}})
	return nil
}

func (acc *Account) RemoveMcAccount(id uuid.UUID) {
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$pull": bson.M{"mc_accounts": id.String()}})
}

func (acc *Account) SetUsername(username string) {
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"username": username}})

	usernameCache.Delete(acc.ID)
}

func (acc *Account) SetEmail(email string) {
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"email": email}})
}

func (acc *Account) SetPassword(password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"password": pass}})
	return nil
}

func (acc *Account) SetCape(id string) {
	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"cape": id}})
}

func (acc *Account) GiveDonator(amount float64) {
	isMember := discord.IsMember(acc.DiscordID)

	username := acc.Username
	if isMember {
		username = fmt.Sprintf("<@%s>", acc.DiscordID)
	}

	discord.SendDonationMessage(username, amount)

	if acc.Donator {
		return
	}

	DonatorCount++

	_, _ = accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$set": bson.M{"donator": true, "can_have_custom_cape": true, "cape": "donator"}})

	if isMember {
		discord.AddRole(acc.DiscordID, discord.AccountRole)
		discord.AddRole(acc.DiscordID, discord.DonorRole)
		discord.SendDonorMsg(acc.DiscordID)
	}
}
