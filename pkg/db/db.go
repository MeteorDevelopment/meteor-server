package db

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"meteor-server/pkg/core"
	"reflect"
)

var client *mongo.Client
var db *mongo.Database

var global *mongo.Collection
var accounts *mongo.Collection
var capes *mongo.Collection
var joinStats *mongo.Collection
var addons *mongo.Collection

func encodeUUID(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	return w.WriteString(v.Interface().(uuid.UUID).String())
}

func decodeUUID(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	str, err := r.ReadString()
	if err != nil {
		return err
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(id))
	return nil
}

func Init() {
	tUUID := reflect.TypeOf(uuid.UUID{})

	registry := bson.NewRegistryBuilder().
		RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(encodeUUID)).
		RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(decodeUUID)).
		Build()

	var err error
	client, err = mongo.NewClient(options.Client().SetRegistry(registry).ApplyURI(core.GetPrivateConfig().MongoDBUrl))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("meteor-bot", options.Database())

	global = db.Collection("global", options.Collection())
	accounts = db.Collection("accounts", options.Collection())
	capes = db.Collection("capes", options.Collection())
	joinStats = db.Collection("join-stats", options.Collection())
	addons = db.Collection("addons", options.Collection())

	initAccounts()
}

func Close() {
	_ = client.Disconnect(nil)
}

func IncrementDownloads() {
	_, _ = global.UpdateOne(nil, bson.M{"id": "Stats"}, bson.M{"$inc": bson.M{"downloads": 1}})
	_, _ = joinStats.UpdateOne(nil, bson.M{"id": core.GetDate()}, bson.M{"$inc": bson.M{"downloads": 1}}, options.Update().SetUpsert(true))
}
