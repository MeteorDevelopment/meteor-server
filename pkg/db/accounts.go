package db

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`

	Donator   bool   `json:"donator"`
	DiscordId string `json:"discordId"`

	MaxMcAccounts int         `json:"maxMcAccounts"`
	McAccounts    []uuid.UUID `json:"mcAccounts"`

	Cape              string `json:"cape"`
	CanHaveCustomCape bool   `json:"canHaveCustomCape"`
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
