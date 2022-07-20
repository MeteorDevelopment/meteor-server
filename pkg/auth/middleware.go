package auth

import (
	"context"
	"meteor-server/pkg/core"
	"net/http"
	"strings"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
			id, err := IsTokenValid(token)

			if err == nil {
				next(w, r.WithContext(context.WithValue(r.Context(), "id", id)))
				return
			}
		}

		core.JsonError(w, "Unauthorized.")
	}
}

func TokenAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == core.GetPrivateConfig().Token {
			next(w, r)
			return
		}

		core.JsonError(w, "Unauthorized.")
	}
}
