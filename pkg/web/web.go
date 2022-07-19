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
	api.InitPayPal()

	r := chi.NewRouter()

	// Setup periodic timers
	t := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			api.ValidateOnlinePlayers()
			<-t.C
		}
	}()

	// Middlewares
	r.Use(middleware.RealIP)

	if core.GetConfig().Debug {
		r.Use(middleware.Logger)
	}

	r.Use(Cors)
	r.Use(middleware.SetHeader("Connection", "close"))
	r.Use(middleware.Recoverer)

	// Static
	r.Handle("/capes/*", http.StripPrefix("/capes", http.FileServer(http.Dir("capes"))))

	// Other
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
			r.Get("/getByUuid", auth.TokenAuth(api.GetAccountByMcUuid))

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
		})

		// /api/discord
		r.Route("/discord", func(r chi.Router) {
			r.Post("/userJoined", auth.TokenAuth(api.DiscordUserJoinedHandler))
			r.Post("/userLeft", auth.TokenAuth(api.DiscordUserLeftHandler))
		})

		// /api/payments
		r.Route("/payments", func(r chi.Router) {
			r.Get("/create", auth.Auth(api.CreateOrderHandler))
			r.Get("/cancel", api.CancelOrderHandler)
			r.Post("/confirm", api.ConfirmOrderHandler)
		})

		// /api/addons
		r.Route("/addons", func(r chi.Router) {
			r.Get("/getById", api.GetAddonById)
			r.Get("/search", api.SearchAddons)
		})
	})

	// Run server
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", core.GetConfig().Port),
		Handler:      r,
		WriteTimeout: 6 * time.Second,
		ReadTimeout:  6 * time.Second,
		IdleTimeout:  6 * time.Second,
	}

	fmt.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
