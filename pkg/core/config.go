package core

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port              int      `json:"port"`
	Debug             bool     `json:"debug"`
	Version           string   `json:"version"`
	DevBuildVersion   string   `json:"dev_build_version"`
	McVersion         string   `json:"mc_version"`
	DevBuildMcVersion string   `json:"dev_build_mc_version"`
	MaxDevBuilds      int      `json:"max_dev_builds"`
	Changelog         []string `json:"changelog"`
}

type PrivateConfig struct {
	MongoDBUrl      string `env:"MONGO_URL"`
	EmailPassword   string `env:"EMAIL_PSW"`
	DiscordToken    string `env:"DISCORD_TOKEN"`
	Token           string `env:"BACKEND_TOKEN"`
	PayPalClientID  string `env:"PAYPAL_CID"`
	PayPalSecret    string `env:"PAYPAL_SECRET"`
	PayPalWebhookId string `env:"PAYPAL_WHID"`
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

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	privateConfig = PrivateConfig{}
	err = env.Parse(&privateConfig)
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
