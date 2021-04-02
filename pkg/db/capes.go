package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Cape struct {
	ID             string `bson:"id"`
	Url            string `bson:"url"`
	SelfAssignable bool   `bson:"selfAssignable"`
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
