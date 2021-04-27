package db

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"meteor-server/pkg/core"
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
		Donator: true,

		DiscordID: "",

		MaxMcAccounts: 1,
		McAccounts:    []uuid.UUID{},

		Cape:              "donator",
		CanHaveCustomCape: true,
	})

	return err
}

func GetAccount(r *http.Request) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"id": core.GetAccountID(r)}).Decode(&acc)

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

func (acc *Account) PasswordMatches(password string) bool {
	return bcrypt.CompareHashAndPassword(acc.Password, []byte(password)) == nil
}

/*func (acc *Account) AddMcAccount(id uuid.UUID) error {
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

	accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$push": bson.M{"mcAccounts": id.String()}})
	return nil
}

func (acc *Account) RemoveMcAccount(id uuid.UUID) {
	accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$pull": bson.M{"mcAccounts": id.String()}})
}*/
