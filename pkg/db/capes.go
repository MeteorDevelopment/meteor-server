package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Cape struct {
	ID             string
	Url            string
	SelfAssignable bool
}

func GetAllCapes() []Cape {
	cursor, err := capes.Find(nil, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var c []Cape
	err = cursor.All(nil, &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
