package clientauth

import (
	"errors"
	"testing"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientAuthenticationSuite struct {
	suite.Suite
}

func TestClientAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(ClientAuthenticationSuite))
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailWhenNoClientIDOrPassKey() {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("", "")
	require.Error(s.T(), err, "should have failed to authenticate")

	assert.Equal(s.T(), "either of Client-ID or Pass-Key is missing, Client-Id: ", err.Error())
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailWhenNoClientID() {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("", "hello")
	require.Error(s.T(), err, "should have failed to authenticate")

	assert.Equal(s.T(), "either of Client-ID or Pass-Key is missing, Client-Id: ", err.Error())
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailWhenNoPassKey() {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("hello", "")
	require.Error(s.T(), err, "should have failed to authenticate")

	assert.Equal(s.T(), "either of Client-ID or Pass-Key is missing, Client-Id: hello", err.Error())
}

func (s *ClientAuthenticationSuite) TestAuthenticateSucceedWhenValidCredsFound() {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	client := &client{ClientID: "hello", PassKey: "bello"}

	clientRepo.On("getClient", "hello").Return(client, nil)

	err := ca.Authenticate("hello", "bello")
	require.NoError(s.T(), err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	require.True(s.T(), found)

	assert.Equal(s.T(), "bello", passKey)

	clientRepo.AssertExpectations(s.T())
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailWhenDBQueryFails() {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	clientRepo.On("getClient", "hello").Return(nil, errors.New("failed to get value from db"))

	err := ca.Authenticate("hello", "bello")
	require.Error(s.T(), err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	require.False(s.T(), found)

	assert.Empty(s.T(), passKey)

	clientRepo.AssertExpectations(s.T())
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailWhenInvalidCreds() {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	client := &client{ClientID: "hello", PassKey: "bello"}

	clientRepo.On("getClient", "hello").Return(client, nil)

	err := ca.Authenticate("hello", "wrong")
	require.Error(s.T(), err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	assert.Equal(s.T(), "bello", passKey)
	require.True(s.T(), found)

	assert.Equal(s.T(), "PassKey is not valid", err.Error())

	clientRepo.AssertExpectations(s.T())
}

func (s *ClientAuthenticationSuite) TestAuthenticateSucceedCacheHitWithValidCreds() {
	myCache := cache.New(cache.NoExpiration, 0)
	myCache.Set("hello", "bello", cache.NoExpiration)

	ca := &ClientAuthentication{
		cache: myCache,
	}

	err := ca.Authenticate("hello", "bello")
	require.NoError(s.T(), err, "should not have failed to authenticate")
}

func (s *ClientAuthenticationSuite) TestAuthenticateFailCacheHitWithInValidCreds() {
	myCache := cache.New(cache.NoExpiration, 0)
	myCache.Set("hello", "bello", cache.NoExpiration)

	ca := &ClientAuthentication{
		cache: myCache,
	}

	err := ca.Authenticate("hello", "wrong")
	require.Error(s.T(), err, "should have failed to authenticate")

	assert.Equal(s.T(), "PassKey is not valid", err.Error())
}
