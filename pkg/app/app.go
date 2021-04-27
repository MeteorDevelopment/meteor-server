package app

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"meteor-server/pkg/auth"
	"net/http"
	"os"
	"time"

	"meteor-server/pkg/api"
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

		return
	}

	if devBuild == "latest" {
		devBuild = db.GetGlobal().DevBuild
	}

	http.Redirect(w, r, fmt.Sprintf("https://%s-309730396-gh.circle-artifacts.com/0/build/libs/meteor-client-%s-%s.jar", devBuild, version, devBuild), http.StatusPermanentRedirect)
}

func Main() {
	core.LoadConfig()
	core.InitEmail()

	db.Init()
	defer db.Close()

	api.UpdateCapes()

	r := mux.NewRouter()

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/favicon.ico", fileHandler("static/assets/favicon.ico"))

	// Redirects
	r.HandleFunc("/discord", redirectHandler("https://discord.com/invite/hv6nz7WScU")).Methods()
	r.HandleFunc("/donate", redirectHandler("https://www.paypal.com/paypalme/MineGame159"))
	r.HandleFunc("/youtube", redirectHandler("https://www.youtube.com/channel/UCWfwmiYGlXXunsUc1Zvz8SQ"))
	r.HandleFunc("/github", redirectHandler("https://github.com/MeteorDevelopment"))

	// Pages
	r.HandleFunc("/", fileHandler("pages/index.html"))
	r.HandleFunc("/info", fileHandler("pages/info.html"))
	r.HandleFunc("/register", fileHandler("pages/register.html"))
	r.HandleFunc("/confirm", fileHandler("pages/confirm.html"))
	r.HandleFunc("/login", fileHandler("pages/login.html"))
	r.HandleFunc("/account", fileHandler("pages/account.html"))

	// Download
	r.HandleFunc("/download", downloadHandler)

	{
		// /api
		g := r.PathPrefix("/api").Subrouter()

		g.HandleFunc("/capes", api.CapesHandler)
		g.HandleFunc("/stats", api.StatsHandler)
		g.HandleFunc("/capeowners", api.CapeOwnersHandler)

		{
			// /api/account
			g2 := g.PathPrefix("/account").Subrouter()

			g2.HandleFunc("/register", api.RegisterHandler).Methods("POST")
			g2.HandleFunc("/confirm", api.ConfirmEmailHandler).Methods("POST")
			g2.HandleFunc("/login", api.LoginHandler)

			g2.HandleFunc("/info", auth.Auth(api.AccountInfoHandler))
			g2.HandleFunc("/mcAccount/:uuid", auth.Auth(api.McAccountHandler)).Methods("POST")
			g2.HandleFunc("/mcAccount/:uuid", auth.Auth(api.McAccountHandler)).Methods("DELETE")
			g2.HandleFunc("/logout", auth.Auth(api.LogoutHandler)).Methods("POST")
		}

		{
			// /api/online
			g2 := g.PathPrefix("/online").Subrouter()

			g2.HandleFunc("/ping", api.PingHandler).Methods("GET", "POST")
			g2.HandleFunc("/leave", api.LeaveHandler).Methods("POST")
			g2.HandleFunc("/usingMeteor", api.UsingMeteorHandler).Methods("POST")
		}
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.GetConfig().Port),
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	fmt.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
