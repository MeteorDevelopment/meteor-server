package db

import "go.mongodb.org/mongo-driver/bson"

type Global struct {
	Downloads      int    `bson:"downloads"`
	TotalAccounts  int    `bson:"totalAccounts"`
	SupportMessage int64  `bson:"supportMessage"`
	DevBuild       string `bson:"devBuild"`
}

func GetGlobal() Global {
	var g Global
	global.FindOne(nil, bson.M{"id": "Stats"}).Decode(&g)
	return g
}
