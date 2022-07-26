package db

import (
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetAddon(id string) (Addon, error) {
	var addon Addon
	err := addons.FindOne(nil, bson.M{"id": id}).Decode(&addon)

	return addon, err
}

func SearchAddons(text string) (*mongo.Cursor, error) {
	if text == "" {
		return addons.Find(nil, bson.M{})
	}

	return addons.Find(nil, bson.M{"title": bson.M{"$regex": text, "$options": "i"}})
}
