package middleware

import (
	"net/http"

	"github.com/KaiserWerk/Maestro/internal/global"
)

func Auth(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Registry-Token")
		if authToken == "" {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}
		if authToken != global.GetAuthToken() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}
