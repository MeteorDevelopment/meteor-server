package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

type JoinStats struct {
	ID        string `json:"date"`
	Joins     int    `json:"joins"`
	Leaves    int    `json:"leaves"`
	Downloads int    `json:"downloads"`
}

func GetJoinStats(date string) (JoinStats, error) {
	var stats JoinStats
	var err = joinStats.FindOne(nil, bson.M{"id": date}).Decode(&stats)

	return stats, err
}
