package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"meteor-server/pkg/core"
)

type JoinStats struct {
	ID        string `bson:"id" json:"date"`
	Joins     int    `bson:"joins" json:"joins"`
	Leaves    int    `bson:"leaves" json:"leaves"`
	Downloads int    `bson:"downloads" json:"downloads"`
}

func GetJoinStats(date string) (JoinStats, error) {
	var stats JoinStats
	var err = joinStats.FindOne(nil, bson.M{"id": date}).Decode(&stats)

	return stats, err
}

func IncrementJoins() {
	_, _ = joinStats.UpdateOne(nil, bson.M{"id": core.GetDate()}, bson.M{"$inc": bson.M{"joins": 1}})
}

func IncrementLeaves() {
	_, _ = joinStats.UpdateOne(nil, bson.M{"id": core.GetDate()}, bson.M{"$inc": bson.M{"leaves": 1}})
}
