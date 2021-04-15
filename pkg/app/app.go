package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"meteor-server/pkg/api"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
)

func fileHandler(file string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(file)
	}
}

func redirectHandler(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, url)
	}
}

func Main() {
	core.LoadConfig()
	core.InitEmail()

	db.Init()
	defer db.Close()

	api.UpdateCapes()

	r := gin.Default()
	r.Static("/static", "static")

	// Redirects
	r.GET("/discord", redirectHandler("https://discord.com/invite/hv6nz7WScU"))
	r.GET("/donate", redirectHandler("https://www.paypal.com/paypalme/MineGame159"))
	r.GET("/youtube", redirectHandler("https://www.youtube.com/channel/UCWfwmiYGlXXunsUc1Zvz8SQ"))
	r.GET("/github", redirectHandler("https://github.com/MeteorDevelopment"))

	// Pages
	r.GET("/", fileHandler("pages/index.html"))
	r.GET("/info", fileHandler("pages/info.html"))
	r.GET("/register", fileHandler("pages/register.html"))
	r.GET("/confirm", fileHandler("pages/confirm.html"))
	r.GET("/login", fileHandler("pages/login.html"))
	r.GET("/account", fileHandler("pages/account.html"))

	{
		// /api
		g := r.Group("/api")

		g.GET("/capes", api.CapesHandler)
		g.GET("/stats", api.StatsHandler)
		g.GET("/capeowners", api.CapeOwnersHandler)

		{
			// /api/account
			g2 := g.Group("/account")

			g2.POST("/register", api.RegisterHandler)
			g2.POST("/confirm", api.ConfirmEmailHandler)
			g2.GET("/login", api.LoginHandler)

			g2.GET("/info", auth.Auth, api.AccountInfoHandler)
			g2.POST("/mcAccount/:uuid", auth.Auth, api.McAccountHandler)
			g2.DELETE("/mcAccount/:uuid", auth.Auth, api.McAccountHandler)
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
