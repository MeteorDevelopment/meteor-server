package db

import (
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	ID       string
	Email    string
	Username string
	Password string

	Donator   bool
	DiscordId string

	MaxMcAccounts int
	McAccounts    []uuid.UUID

	Cape              string
	CanHaveCustomCape bool
}

func GetAccountsWithCape() []Account {
	cursor, err := accounts.Find(nil, bson.M{"cape": bson.M{"$ne": ""}})
	if err != nil {
		log.Fatal(err)
	}

	var a []Account
	err = cursor.All(nil, &a)
	if err != nil {
		log.Fatal(err)
	}

	return a
}

func GetAccountWithUsernameOrEmail(name string) (Account, error) {
	var a Account
	err := accounts.FindOne(nil, bson.M{"username": name}).Decode(&a)

	if err != nil {
		err = accounts.FindOne(nil, bson.M{"email": name}).Decode(&a)
	}

	return a, err
}
