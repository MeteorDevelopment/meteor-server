package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"
)

type Global struct {
	Downloads       int    `bson:"downloads"`
	TotalAccounts   int    `bson:"totalAccounts"`
	SupportMessage  int64  `bson:"supportMessage"`
	DevBuild        string `bson:"devBuild"`
	DevBuildVersion string `bson:"devBuildVersion"`
}

var cache Global
var lastTime time.Time

func GetGlobal() Global {
	now := time.Now()

	if now.Sub(lastTime) > time.Second {
		err := global.FindOne(nil, bson.M{"id": "Stats"}).Decode(&cache)
		lastTime = now

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to query global stats:", err)
		}
	}

	return cache
}

func SetDevBuild(devBuild string) {
	_, _ = global.UpdateOne(nil, bson.M{"id": "Stats"}, bson.M{"$set": bson.M{"devBuild": devBuild}})
}

func SetDevBuildVersion(version string) {
	_, _ = global.UpdateOne(nil, bson.M{"id": "Stats"}, bson.M{"$set": bson.M{"devBuildVersion": version}})
}
