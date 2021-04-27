package main

import (
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"meteor-server/pkg/web"
)

func main() {
	core.LoadConfig()
	core.InitEmail()

	db.Init()
	defer db.Close()

	web.Main()
}
