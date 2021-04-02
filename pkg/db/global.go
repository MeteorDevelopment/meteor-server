package db

import "go.mongodb.org/mongo-driver/bson"

type Global struct {
	Downloads      int
	TotalAccounts  int
	SupportMessage int64
	DevBuild       string
}

func GetGlobal() Global {
	var g Global
	global.FindOne(nil, bson.M{"id": "Stats"}).Decode(&g)
	return g
}
