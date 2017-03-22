package clientIdAuth

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
)

type middleware func(http.Handler) http.Handler

var RequestIdCache = cache.New(cache.NoExpiration, 0)

func BuildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

func WithClientIdAndPassKeyAuthorization(clientRepository ClientRepositoryInterface) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewV4().String()
			context.Set(r, "requestID", requestID)
			logFields := BuildContext("authMiddleware")
			requestClientId := r.Header.Get("Client-ID")
			requestPassKey := r.Header.Get("Pass-Key")
			cachedPassKey, _ := RequestIdCache.Get(requestClientId)

			if requestClientId == "" || requestPassKey == "" {
				logrus.WithFields(logFields).Error("Either of client-id or paskey is missing, client-id: "+requestClientId, " pass-key: ", requestPassKey)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if cachedPassKey == "" {
				authorised_application, dbError := clientRepository.GetClient(requestClientId)
				if dbError != nil {
					logrus.WithFields(logFields).Error("Error fetching client ID from DB" + dbError.Error())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				RequestIdCache.Set(authorised_application.ClientId, authorised_application.PassKey, cache.NoExpiration)
				cachedPassKey = authorised_application.PassKey
			}

			if cachedPassKey != requestPassKey {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
