package clientauth

import (
	"errors"
	"testing"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticateFailWhenNoClientIDOrPassKey(t *testing.T) {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("", "")
	require.Error(t, err, "should have failed to authenticate")

	assert.Equal(t, "either of Client-ID or Pass-Key is missing, Client-Id: ", err.Error())
}

func TestAuthenticateFailWhenNoClientID(t *testing.T) {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("", "hello")
	require.Error(t, err, "should have failed to authenticate")

	assert.Equal(t, "either of Client-ID or Pass-Key is missing, Client-Id: ", err.Error())
}

func TestAuthenticateFailWhenNoPassKey(t *testing.T) {
	authConfig := NewConfig("postgres", "dbname=client_auth sslmode=disable")
	ca := NewClientAuthentication(authConfig)

	err := ca.Authenticate("hello", "")
	require.Error(t, err, "should have failed to authenticate")

	assert.Equal(t, "either of Client-ID or Pass-Key is missing, Client-Id: hello", err.Error())
}

func TestAuthenticateSucceedWhenValidCredsFound(t *testing.T) {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	client := &client{ClientID: "hello", PassKey: "bello"}

	clientRepo.On("getClient", "hello").Return(client, nil)

	err := ca.Authenticate("hello", "bello")
	require.NoError(t, err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	require.True(t, found)

	assert.Equal(t, "bello", passKey)

	clientRepo.AssertExpectations(t)
}

func TestAuthenticateFailWhenDBQueryFails(t *testing.T) {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	clientRepo.On("getClient", "hello").Return(nil, errors.New("failed to get value from db"))

	err := ca.Authenticate("hello", "bello")
	require.Error(t, err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	require.False(t, found)

	assert.Empty(t, passKey)

	clientRepo.AssertExpectations(t)
}

func TestAuthenticateFailWhenInvalidCreds(t *testing.T) {
	clientRepo := &mockClientRepository{}
	myCache := cache.New(cache.NoExpiration, 0)

	ca := &ClientAuthentication{
		cache: myCache,
		db:    clientRepo,
	}

	client := &client{ClientID: "hello", PassKey: "bello"}

	clientRepo.On("getClient", "hello").Return(client, nil)

	err := ca.Authenticate("hello", "wrong")
	require.Error(t, err, "should have failed to authenticate")

	passKey, found := myCache.Get("hello")
	assert.Equal(t, "bello", passKey)
	require.True(t, found)

	assert.Equal(t, "PassKey is not valid", err.Error())

	clientRepo.AssertExpectations(t)
}

func TestAuthenticateSucceedCacheHitWithValidCreds(t *testing.T) {
	myCache := cache.New(cache.NoExpiration, 0)
	myCache.Set("hello", "bello", cache.NoExpiration)

	ca := &ClientAuthentication{
		cache: myCache,
	}

	err := ca.Authenticate("hello", "bello")
	require.NoError(t, err, "should not have failed to authenticate")
}

func TestAuthenticateFailCacheHitWithInValidCreds(t *testing.T) {
	myCache := cache.New(cache.NoExpiration, 0)
	myCache.Set("hello", "bello", cache.NoExpiration)

	ca := &ClientAuthentication{
		cache: myCache,
	}

	err := ca.Authenticate("hello", "wrong")
	require.Error(t, err, "should have failed to authenticate")

	assert.Equal(t, "PassKey is not valid", err.Error())
}
