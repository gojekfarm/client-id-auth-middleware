package clientauth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var called bool

type nextHandler struct{}

func (nh nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	called = true
}

func TestAuthorizationWhenClientIDMissing(t *testing.T) {
	r, err := http.NewRequest("GET", "/authenticate", nil)
	require.NoError(t, err, "should not have failed to create a request")

	r.Header.Add("Pass-Key", "some key")

	w := httptest.NewRecorder()

	mockClientAuthenticator := &mockClientAuthenticator{}
	mockClientAuthenticator.On("Authenticate", "", "some key").Return(errors.New("failed to authorize client"))

	WithClientIDAndPassKeyAuthorization(mockClientAuthenticator)(nextHandler{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, called)

	mockClientAuthenticator.AssertExpectations(t)
}

func TestAuthorizationWhenPassKeyMissing(t *testing.T) {
	r, err := http.NewRequest("GET", "/authenticate", nil)
	require.NoError(t, err, "should not have failed to create a request")

	r.Header.Add("Client-ID", "some_client_id")

	w := httptest.NewRecorder()

	mockClientAuthenticator := &mockClientAuthenticator{}
	mockClientAuthenticator.On("Authenticate", "some_client_id", "").Return(errors.New("failed to authorize client"))

	WithClientIDAndPassKeyAuthorization(mockClientAuthenticator)(nextHandler{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, called)

	mockClientAuthenticator.AssertExpectations(t)
}

func TestAuthorizationFailWithInvalidCreds(t *testing.T) {
	r, err := http.NewRequest("GET", "/authenticate", nil)
	require.NoError(t, err, "should not have failed to create a request")

	r.Header.Add("Client-ID", "some_client_id")
	r.Header.Add("Pass-Key", "some_pass_key")

	w := httptest.NewRecorder()

	mockClientAuthenticator := &mockClientAuthenticator{}
	mockClientAuthenticator.On("Authenticate", "some_client_id", "some_pass_key").Return(errors.New("failed to authorize client"))

	WithClientIDAndPassKeyAuthorization(mockClientAuthenticator)(nextHandler{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, called)

	mockClientAuthenticator.AssertExpectations(t)
}

func TestAuthorizationSucceed(t *testing.T) {
	r, err := http.NewRequest("GET", "/authenticate", nil)
	require.NoError(t, err, "should not have failed to create a request")

	r.Header.Add("Client-ID", "some_client_id")
	r.Header.Add("Pass-Key", "some_pass_key")

	w := httptest.NewRecorder()

	mockClientAuthenticator := &mockClientAuthenticator{}
	mockClientAuthenticator.On("Authenticate", "some_client_id", "some_pass_key").Return(nil)

	WithClientIDAndPassKeyAuthorization(mockClientAuthenticator)(nextHandler{}).ServeHTTP(w, r)

	assert.True(t, called)

	mockClientAuthenticator.AssertExpectations(t)
}
