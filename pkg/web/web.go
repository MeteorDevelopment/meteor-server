package web

import (
	"fmt"
	"io/fs"
	"log"
	"meteor-server/pkg/auth"
	"meteor-server/pkg/web/api"
	"meteor-server/pkg/wormhole"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"meteor-server/pkg/core"
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

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func Main() {
	err := os.Mkdir("capes", fs.ModeDir)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	api.UpdateCapes()

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RealIP)

	if core.GetConfig().Debug {
		r.Use(middleware.Logger)
	}

	r.Use(Cors)
	r.Use(middleware.SetHeader("Connection", "close"))
	r.Use(middleware.Recoverer)

	// Static
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	r.Handle("/capes/*", http.StripPrefix("/capes", http.FileServer(http.Dir("capes"))))

	// Redirects
	r.Get("/discord", redirectHandler("https://discord.com/invite/bBGQZvd"))
	r.Get("/donate", redirectHandler("https://www.paypal.com/paypalme/MineGame159"))
	r.Get("/youtube", redirectHandler("https://www.youtube.com/channel/UCWfwmiYGlXXunsUc1Zvz8SQ"))
	r.Get("/github", redirectHandler("https://github.com/MeteorDevelopment"))
	r.Get("/faq", redirectHandler("https://github.com/MeteorDevelopment/meteor-client/wiki"))

	// Pages
	r.Get("/", fileHandler("pages/index.html"))
	r.Get("/changelog", fileHandler("pages/changelog.html"))
	r.Get("/donations", fileHandler("pages/donations.html"))
	r.Get("/register", fileHandler("pages/register.html"))
	r.Get("/confirm", fileHandler("pages/confirm.html"))
	r.Get("/login", fileHandler("pages/login.html"))
	r.Get("/account", fileHandler("pages/account.html"))
	r.Get("/changeUsername", fileHandler("pages/changeUsername.html"))
	r.Get("/changeEmail", fileHandler("pages/changeEmail.html"))
	r.Get("/confirmChangeEmail", api.ConfirmChangeEmailHandler)
	r.Get("/changePassword", fileHandler("pages/changePassword.html"))
	r.Get("/forgotPassword", fileHandler("pages/forgotPassword.html"))

	// Other
	r.Get("/favicon.ico", fileHandler("static/assets/favicon.ico"))
	r.Get("/icon.png", fileHandler("static/assets/icon.png"))
	r.Get("/download", api.DownloadHandler)

	if core.GetConfig().Debug {
		r.Get("/handler.go", wormhole.Handle)
	}

	// /api
	r.Route("/api", func(r chi.Router) {
		r.Get("/download", api.DownloadHandler)

		r.Get("/capes", api.CapesHandler)
		r.Get("/stats", api.StatsHandler)
		r.Get("/capeowners", api.CapeOwnersHandler)

		r.Post("/uploadDevBuild", auth.TokenAuth(api.UploadDevBuildHandler))

		// /api/account
		r.Route("/account", func(r chi.Router) {
			r.Post("/register", api.RegisterHandler)
			r.Post("/confirm", api.ConfirmEmailHandler)
			r.Get("/login", api.LoginHandler)
			r.Post("/forgotPassword", api.ForgotPasswordHandler)
			r.Post("/logout", auth.Auth(api.LogoutHandler))

			r.Get("/info", auth.Auth(api.AccountInfoHandler))
			r.Get("/generateDiscordLinkToken", auth.Auth(api.GenerateDiscordLinkTokenHandler))
			r.Post("/linkDiscord", auth.TokenAuth(api.LinkDiscordHandler))
			r.Post("/unlinkDiscord", auth.Auth(api.UnlinkDiscordHandler))
			r.Post("/mcAccount", auth.Auth(api.McAccountHandler))
			r.Delete("/mcAccount", auth.Auth(api.McAccountHandler))
			r.Post("/selectCape", auth.Auth(api.SelectCapeHandler))
			r.Post("/uploadCape", auth.Auth(api.UploadCapeHandler))
			r.Post("/changeUsername", auth.Auth(api.ChangeUsernameHandler))
			r.Post("/changeEmail", auth.Auth(api.ChangeEmailHandler))
			r.Get("/confirmChangeEmail", api.ConfirmChangeEmailHandlerApi)
			r.Post("/changePassword", auth.Auth(api.ChangePasswordHandler))
			r.Post("/changePasswordToken", api.ChangePasswordTokenHandler)
		})

		// /api/online
		r.Route("/online", func(r chi.Router) {
			r.Get("/ping", api.PingHandler)
			r.Post("/ping", api.PingHandler)
			r.Post("/leave", api.LeaveHandler)
			r.Post("/usingMeteor", api.UsingMeteorHandler)
		})

		// /api/discord
		r.Route("/discord", func(r chi.Router) {
			r.Post("/userJoined", auth.TokenAuth(api.DiscordUserJoinedHandler))
			r.Post("/userLeft", auth.TokenAuth(api.DiscordUserLeftHandler))
			r.Post("/giveDonator", auth.TokenAuth(api.GiveDonatorHandler))
		})
	})

	// Run server
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.GetConfig().Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	fmt.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
