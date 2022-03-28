package core

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port            int      `json:"port"`
	Debug           bool     `json:"debug"`
	Version         string   `json:"version"`
	DevBuildVersion string   `json:"dev_build_version"`
	McVersion       string   `json:"mc_version"`
	MaxDevBuilds    int      `json:"max_dev_builds"`
	Changelog       []string `json:"changelog"`
}

type PrivateConfig struct {
	MongoDBUrl    string
	EmailPassword string
	DiscordToken  string
	Token         string
}

var config Config
var privateConfig PrivateConfig

func LoadConfig() {
	// Config
	f, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Private config
	f, err = os.ReadFile("private_config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &privateConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() Config {
	return config
}

func GetPrivateConfig() PrivateConfig {
	return privateConfig
}
