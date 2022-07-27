package db

import (
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"time"
)

type Addon struct {
	ID string `bson:"id" json:"id"`

	Title       string `bson:"title" json:"title"`
	Icon        string `bson:"icon" json:"icon"`
	Description string `bson:"description" json:"description"`
	Markdown    string `bson:"markdown" json:"markdown"`

	Developers []ksuid.KSUID `bson:"developers" json:"developers"`

	Version        string   `bson:"version" json:"version"`
	MeteorVersions []string `bson:"meteor_versions" json:"meteor_versions"`
	Download       string   `bson:"download" json:"download"`

	DownloadCount int `bson:"download_count" json:"download_count"`

	Website string `bson:"website" json:"website"`
	Source  string `bson:"source" json:"source"`
	Support string `bson:"support" json:"support"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

const pageSize = 15

func GetAddon(id string) (Addon, error) {
	var addon Addon
	err := addons.FindOne(nil, bson.M{"id": id}).Decode(&addon)

	return addon, err
}

func SearchAddons(text string, page int) (*mongo.Cursor, int, error) {
	if page < 1 {
		page = 1
	}

	sort := bson.D{{"$sort", bson.D{
		{"download_count", -1},
	}}}
	skip := bson.D{{"$skip", (page - 1) * pageSize}}
	limit := bson.D{{"$limit", pageSize}}

	if text == "" {
		count, err := addons.EstimatedDocumentCount(nil)
		if err != nil {
			return nil, 0, err
		}

		cursor, err := addons.Aggregate(nil, mongo.Pipeline{
			sort,
			skip,
			limit,
		})

		return cursor, int(math.Ceil(float64(count) / pageSize)), err
	}

	count, err := addons.CountDocuments(nil, bson.M{"title": bson.M{"$regex": text, "$options": "i"}})
	if err != nil {
		return nil, 0, err
	}

	cursor, err := addons.Aggregate(nil, mongo.Pipeline{
		bson.D{{"$match", bson.D{
			{"title", bson.M{
				"$regex":   text,
				"$options": "i",
			}},
		}}},
		sort,
		skip,
		limit,
	})

	return cursor, int(math.Ceil(float64(count) / pageSize)), err
}
