package middleware

import (
	"github.com/KaiserWerk/Maestro/internal/configuration"
	"net/http"
)

func Auth(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Registry-Token")
		if authToken == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}
		if conf := configuration.GetConfiguration(); authToken != conf.App.AuthToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}
