package middleware

import (
	"errors"
	"net/http"

	"github.com/caesariodito/caesariodito.dev-backend/api"
	"github.com/caesariodito/caesariodito.dev-backend/internal/tools"
	log "github.com/sirupsen/logrus"
)

var errUnauthorized = errors.New("invalid username or token")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token string = r.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(errUnauthorized)
			api.RequestErrorHandler(w, errUnauthorized)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		loginDetails := (*database).GetUserLoginDetails(username)

		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(errUnauthorized)
			api.RequestErrorHandler(w, errUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
