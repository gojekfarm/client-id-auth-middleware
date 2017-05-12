package clientauth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
)

type middleware func(http.Handler) http.Handler

var requestIdCache = cache.New(cache.NoExpiration, 0)

func buildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

func WithClientIdAndPassKeyAuthorization(clientRepository ClientRepositoryInterface) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewV4().String()
			context.Set(r, "requestID", requestID)
			logFields := buildContext("authMiddleware")
			requestClientID := r.Header.Get("Client-ID")
			requestPassKey := r.Header.Get("Pass-Key")
			cachedPassKey, _ := requestIdCache.Get(requestClientID)

			if requestClientID == "" || requestPassKey == "" {
				logrus.WithFields(logFields).Error("Either of client-id or paskey is missing, client-id: "+requestClientID, " pass-key: ", requestPassKey)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if cachedPassKey == "" {
				authorizedApplication, dbError := clientRepository.GetClient(requestClientID)
				if dbError != nil {
					logrus.WithFields(logFields).Error("Error fetching client ID from DB" + dbError.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				requestIdCache.Set(authorizedApplication.ClientId, authorizedApplication.PassKey, cache.NoExpiration)
				cachedPassKey = authorizedApplication.PassKey
			}

			if cachedPassKey != requestPassKey {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
