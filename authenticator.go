package clientauth

import (
	"errors"
	"fmt"

	cache "github.com/patrickmn/go-cache"
)

type ClientAuthenticator interface {
	Authenticate(clientID, passKey string) error
	HeaderConfig() *HeaderConfig
}

type ClientAuthentication struct {
	cache  *cache.Cache
	db     clientStore
	config *Config
}

func NewClientAuthentication(authConfig *Config) *ClientAuthentication {
	return &ClientAuthentication{
		cache: cache.New(cache.NoExpiration, 0),
		db: &clientRepository{
			db: loadDatabase(authConfig.dbDriver, authConfig.dbConnURL),
		},
		config: authConfig,
	}
}

func (ca *ClientAuthentication) HeaderConfig() *HeaderConfig {
	return ca.config.HeaderConfig
}

func (ca *ClientAuthentication) Authenticate(clientID, passKey string) error {
	if len(clientID) <= 0 || len(passKey) <= 0 {
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
