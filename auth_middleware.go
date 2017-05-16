package clientauth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/satori/go.uuid"
)

type Middleware func(http.Handler) http.Handler
type NextMiddleware func(http.ResponseWriter, *http.Request, http.Handler)

func WithClientIDAndPassKeyAuthorization(authenticator ClientAuthenticator) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//TODO: Use golang context for setting request specific data
			requestID := uuid.NewV4().String()
			context.Set(r, "requestID", requestID)

			//TODO: Take in the logger from client as a config
			logger := logrus.WithFields(buildContext("authMiddleware"))

			//TODO: Make clientid and passkey headers as configurable
			err := authenticator.Authenticate(readAuthHeaders(r.Header))
			if err != nil {
				logger.Errorf("failed to authenticate client")

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func NextAuthorizer(authenticator ClientAuthenticator) NextMiddleware {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		err := authenticator.Authenticate(readAuthHeaders(r.Header))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func readAuthHeaders(headers http.Header) (string, string) {
	return headers.Get("Client-ID"), headers.Get("Pass-Key")
}

func buildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}
