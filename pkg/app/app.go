package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"meteor-server/pkg/api"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
)

func indexHandler(c *gin.Context) {
	c.File("pages/index.html")
}

func infoHandler(c *gin.Context) {
	c.File("pages/info.html")
}

func Main() {
	core.LoadConfig()
	core.InitEmail()

	db.Init()
	defer db.Close()

	api.UpdateCapes()

	r := gin.Default()
	r.Static("/static", "static")

	r.GET("/", indexHandler)
	r.GET("/info", infoHandler)

	{
		// /api
		g := r.Group("/api")

		g.GET("/capes", api.CapesHandler)
		g.GET("/stats", api.StatsHandler)
		g.GET("/capeowners", api.CapeOwnersHandler)

		{
			// /api/account
			g2 := g.Group("/account")

			g2.GET("/login", api.LoginHandler)

			g2.POST("/logout", auth.Auth, api.LogoutHandler)
		}

		{
			// /api/online
			g2 := g.Group("/online")

			g2.GET("/ping", api.PingHandler) // TODO: Deprecated
			g2.POST("/ping", api.PingHandler)
			g2.POST("/leave", api.LeaveHandler)
			g2.POST("/usingMeteor", api.UsingMeteorHandler)
		}
	}

	log.Fatal(r.Run())
}
