package middleware

import (
	"net/http"

	"github.com/KaiserWerk/Maestro/internal/configuration"
	"github.com/KaiserWerk/Maestro/internal/global"
)

func Auth(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(global.AuthHeader)
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
