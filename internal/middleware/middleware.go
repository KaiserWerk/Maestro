package middleware

import (
	"net/http"

	"github.com/KaiserWerk/Maestro/internal/global"
)

func (mwh *MWHandler) Auth(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	l := mwh.Logger.WithField("context", "middleware.Auth")
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(global.AuthHeader)
		if authToken == "" {
			l.Warn("missing auth token")
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}
		if authToken != mwh.Config.AuthToken {
			l.Warn("failed authentication attempt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}
