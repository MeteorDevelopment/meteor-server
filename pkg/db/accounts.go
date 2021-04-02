package db

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	ID       string `bson:"id" json:"id"`
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"-"`

	Donator   bool   `bson:"donator" json:"donator"`
	DiscordId string `bson:"discordId" json:"discordId"`

	MaxMcAccounts int         `bson:"maxMcAccounts" json:"maxMcAccounts"`
	McAccounts    []uuid.UUID `bson:"mcAccounts" json:"mcAccounts"`

	Cape              string `bson:"cape" json:"cape"`
	CanHaveCustomCape bool   `bson:"canHaveCustomCape" json:"canHaveCustomCape"`
}

func GetAccount(c *gin.Context) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"id": c.GetString("id")}).Decode(&acc)

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

func GetAccountWithUsernameOrEmail(name string) (Account, error) {
	var acc Account
	err := accounts.FindOne(nil, bson.M{"username": name}).Decode(&acc)

	if err != nil {
		err = accounts.FindOne(nil, bson.M{"email": name}).Decode(&acc)
	}

	return acc, err
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

	accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$push": bson.M{"mcAccounts": id.String()}})
	return nil
}

func (acc *Account) RemoveMcAccount(id uuid.UUID) {
	accounts.UpdateOne(nil, bson.M{"id": acc.ID}, bson.M{"$pull": bson.M{"mcAccounts": id.String()}})
}
