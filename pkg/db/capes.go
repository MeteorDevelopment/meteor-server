package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Cape struct {
	ID  string `bson:"id" json:"id"`
	Url string `bson:"url" json:"url"`
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

func GetCape(id string) (Cape, error) {
	var cape Cape
	err := capes.FindOne(nil, bson.M{"id": id}).Decode(&cape)

	return cape, err
}

func InsertCape(cape Cape) {
	_, _ = capes.DeleteOne(nil, bson.M{"id": cape.ID})
	_, _ = capes.InsertOne(nil, cape)
}
