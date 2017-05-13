package clientauth

import (
	"errors"
	"fmt"

	cache "github.com/patrickmn/go-cache"
)

type clientAuthenticator interface {
	Authenticate(clientID, passKey string) error
}

type ClientAuthentication struct {
	cache *cache.Cache
	db    clientStore
}

func NewClientAuthentication(authConfig *Config) *ClientAuthentication {
	return &ClientAuthentication{
		cache: cache.New(cache.NoExpiration, 0),
		db: &clientRepository{
			db: loadDatabase(authConfig.dbDriver, authConfig.dbConnURL),
		},
	}
}

func (ca *ClientAuthentication) Authenticate(clientID, passKey string) error {
	if clientID == "" || passKey == "" {
		return fmt.Errorf("either of Client-ID or Pass-Key is missing, Client-Id: %s", clientID)
	}

	cachedPassKey, found := ca.cache.Get(clientID)
	if !found {
		authorizedApplication, err := ca.db.getClient(clientID)
		if err != nil {
			return fmt.Errorf("failed to query the database: %s", err)
		}

		ca.cache.Set(authorizedApplication.ClientID, authorizedApplication.PassKey, cache.NoExpiration)
		cachedPassKey = authorizedApplication.PassKey
	}

	if cachedPassKey != passKey {
		return errors.New("PassKey is not valid")
	}

	return nil
}
