package clientauth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
)

type middleware func(http.Handler) http.Handler

var requestIDCache = cache.New(cache.NoExpiration, 0)

func buildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

func WithClientIDAndPassKeyAuthorization(clientRepo clientStore) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewV4().String()
			context.Set(r, "requestID", requestID)

			logger := logrus.WithFields(buildContext("authMiddleware"))

			requestClientID := r.Header.Get("Client-ID")
			requestPassKey := r.Header.Get("Pass-Key")

			if requestClientID == "" || requestPassKey == "" {
				logger.Errorf("either of Client-ID or Pass-Key is missing, Client-Id: %s Pass-Key: %s", requestClientID, requestPassKey)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			cachedPassKey, found := requestIDCache.Get(requestClientID)
			if !found {

				authorizedApplication, err := clientRepo.GetClient(requestClientID)
				if err != nil {
					logger.Errorf("error fetching client ID from DB %s", err)

					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				requestIDCache.Set(authorizedApplication.ClientID, authorizedApplication.PassKey, cache.NoExpiration)
				cachedPassKey = authorizedApplication.PassKey
			}

			if cachedPassKey != requestPassKey {
				logger.Error("error: Pass-Key is invalid")

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
