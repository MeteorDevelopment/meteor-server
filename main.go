package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"meteor-server/pkg/web"
	"os"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.TimeOnly,
	})

	core.Init()
	core.LoadConfig()
	core.InitEmail()

	db.Init()
	defer db.Close()

	auth.Init()
	web.Main()
}
