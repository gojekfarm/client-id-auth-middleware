package clientauth

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type Middleware func(http.Handler) http.Handler
type NextMiddleware func(http.ResponseWriter, *http.Request, http.HandlerFunc)

type requestIDKey string

var RequestID requestIDKey = "Request-ID"

func WithClientIDAndPassKeyAuthorization(authenticator ClientAuthenticator) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//TODO: Take in the logger from client as a config
			logger := logrus.WithFields(buildContext("authMiddleware"))

			err := authenticator.Authenticate(readAuthHeaders(r.Header, authenticator.HeaderConfig()))
			if err != nil {
				logger.Errorf("failed to authenticate client: %s", err)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r.WithContext(contextWithRequestID(r)))
		})
	}
}

func NextAuthorizer(authenticator ClientAuthenticator) NextMiddleware {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err := authenticator.Authenticate(readAuthHeaders(r.Header, authenticator.HeaderConfig()))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(contextWithRequestID(r)))
	}
}

func readAuthHeaders(headers http.Header, headerConfig *HeaderConfig) (string, string) {
	return headers.Get(headerConfig.ClientIDName), headers.Get(headerConfig.PassKeyName)
}

func buildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

func contextWithRequestID(r *http.Request) context.Context {
	requestID := uuid.NewV4().String()
	return context.WithValue(r.Context(), RequestID, requestID)
}
