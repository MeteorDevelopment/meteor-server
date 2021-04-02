package db

import (
	"go.mongodb.org/mongo-driver/bson"
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
