package clientauth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/satori/go.uuid"
)

type middleware func(http.Handler) http.Handler

func buildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

func WithClientIDAndPassKeyAuthorization(authenticator ClientAuthenticator) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewV4().String()
			context.Set(r, "requestID", requestID)

			logger := logrus.WithFields(buildContext("authMiddleware"))

			requestClientID := r.Header.Get("Client-ID")
			requestPassKey := r.Header.Get("Pass-Key")

			err := authenticator.Authenticate(requestClientID, requestPassKey)
			if err != nil {
				logger.Errorf("failed to authenticate client for ID : %s", requestClientID)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
