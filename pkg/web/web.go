package web

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/fs"
	"log"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/web/api"
	"meteor-server/pkg/wormhole"
	"net/http"
	"os"
	"time"

	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
)

func fileHandler(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}

func redirectHandler(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	version := core.GetConfig().Version
	devBuild := r.URL.Query().Get("devBuild")

	if devBuild == "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=meteor-client-%s.jar", version))
		http.ServeFile(w, r, fmt.Sprintf("jars/meteor-client-%s.jar", version))

		db.IncrementDownloads()
		return
	}

	version = core.GetConfig().DevBuildVersion

	if devBuild == "latest" {
		devBuild = db.GetGlobal().DevBuild
	}

	http.Redirect(w, r, fmt.Sprintf("https://%s-309730396-gh.circle-artifacts.com/0/build/libs/meteor-client-%s-%s.jar", devBuild, version, devBuild), http.StatusPermanentRedirect)
	db.IncrementDownloads()
}

func Main() {
	err := os.Mkdir("capes", fs.ModeDir)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	api.UpdateCapes()

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			next.ServeHTTP(w, r)
		})
	})

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/capes").Handler(http.StripPrefix("/capes", http.FileServer(http.Dir("capes"))))

	// Redirects
	r.HandleFunc("/discord", redirectHandler("https://discord.com/invite/hv6nz7WScU"))
	r.HandleFunc("/donate", redirectHandler("https://www.paypal.com/paypalme/MineGame159"))
	r.HandleFunc("/youtube", redirectHandler("https://www.youtube.com/channel/UCWfwmiYGlXXunsUc1Zvz8SQ"))
	r.HandleFunc("/github", redirectHandler("https://github.com/MeteorDevelopment"))
	r.HandleFunc("/faq", redirectHandler("https://github.com/MeteorDevelopment/meteor-client/wiki"))

	// Pages
	r.HandleFunc("/", fileHandler("pages/index.html"))
	r.HandleFunc("/changelog", fileHandler("pages/changelog.html"))
	r.HandleFunc("/donations", fileHandler("pages/donations.html"))
	r.HandleFunc("/register", fileHandler("pages/register.html"))
	r.HandleFunc("/confirm", fileHandler("pages/confirm.html"))
	r.HandleFunc("/login", fileHandler("pages/login.html"))
	r.HandleFunc("/account", fileHandler("pages/account.html"))
	r.HandleFunc("/changeUsername", fileHandler("pages/changeUsername.html"))
	r.HandleFunc("/changeEmail", fileHandler("pages/changeEmail.html"))
	r.HandleFunc("/confirmChangeEmail", api.ConfirmChangeEmailHandler)
	r.HandleFunc("/changePassword", fileHandler("pages/changePassword.html"))

	// Other
	r.HandleFunc("/favicon.ico", fileHandler("static/assets/favicon.ico"))
	r.HandleFunc("/download", downloadHandler)
	r.HandleFunc("/handler.go", wormhole.Handle)

	{
		// /api
		g := r.PathPrefix("/api").Subrouter()

		g.HandleFunc("/capes", api.CapesHandler)
		g.HandleFunc("/stats", api.StatsHandler)
		g.HandleFunc("/capeowners", api.CapeOwnersHandler)
		g.HandleFunc("/setDevBuild", auth.TokenAuth(api.SetDevBuildHandler)).Methods("POST")

		{
			// /api/account
			g2 := g.PathPrefix("/account").Subrouter()

			g2.HandleFunc("/register", api.RegisterHandler).Methods("POST")
			g2.HandleFunc("/confirm", api.ConfirmEmailHandler).Methods("POST")
			g2.HandleFunc("/login", api.LoginHandler)
			g2.HandleFunc("/logout", auth.Auth(api.LogoutHandler)).Methods("POST")

			g2.HandleFunc("/info", auth.Auth(api.AccountInfoHandler))
			g2.HandleFunc("/generateDiscordLinkToken", auth.Auth(api.GenerateDiscordLinkTokenHandler))
			g2.HandleFunc("/linkDiscord", auth.TokenAuth(api.LinkDiscordHandler)).Methods("POST")
			g2.HandleFunc("/unlinkDiscord", auth.Auth(api.UnlinkDiscordHandler)).Methods("POST")
			g2.HandleFunc("/mcAccount", auth.Auth(api.McAccountHandler)).Methods("POST", "DELETE")
			g2.HandleFunc("/selectCape", auth.Auth(api.SelectCapeHandler)).Methods("POST")
			g2.HandleFunc("/uploadCape", auth.Auth(api.UploadCapeHandler)).Methods("POST")
			g2.HandleFunc("/changeUsername", auth.Auth(api.ChangeUsernameHandler)).Methods("POST")
			g2.HandleFunc("/changeEmail", auth.Auth(api.ChangeEmailHandler)).Methods("POST")
			g2.HandleFunc("/changePassword", auth.Auth(api.ChangePasswordHandler)).Methods("POST")
		}

		{
			// /api/online
			g2 := g.PathPrefix("/online").Subrouter()

			g2.HandleFunc("/ping", api.PingHandler).Methods("GET", "POST")
			g2.HandleFunc("/leave", api.LeaveHandler).Methods("POST")
			g2.HandleFunc("/usingMeteor", api.UsingMeteorHandler).Methods("POST")
		}

		{
			// /api/discord
			g2 := g.PathPrefix("/discord").Subrouter()

			g2.HandleFunc("/userJoined", auth.TokenAuth(api.DiscordUserJoinedHandler)).Methods("POST")
			g2.HandleFunc("/userLeft", auth.TokenAuth(api.DiscordUserLeftHandler)).Methods("POST")
			g2.HandleFunc("/giveDonator", auth.TokenAuth(api.GiveDonatorHandler)).Methods("POST")
		}
	}

	var handler http.Handler = r
	if core.GetConfig().Debug {
		handler = handlers.LoggingHandler(os.Stdout, handler)
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.GetConfig().Port),
		Handler:      handler,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	fmt.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
