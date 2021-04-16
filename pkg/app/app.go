package app

import (
	"fmt"
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

func downloadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := core.GetConfig().Version
		devBuild := c.Request.URL.Query().Get("devBuild")

		if devBuild == "" {
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=meteor-client-%s.jar", version))
			c.Writer.Header().Set("Content-Type", "application/java-archive")
			c.File(fmt.Sprintf("jars/meteor-client-%s.jar", version))
			return
		}

		if devBuild == "latest" {
			devBuild = db.GetGlobal().DevBuild
		}

		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("https://%s-309730396-gh.circle-artifacts.com/0/build/libs/meteor-client-%s-%s.jar", devBuild, version, devBuild))
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

	// Download
	r.GET("/download", downloadHandler())

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
