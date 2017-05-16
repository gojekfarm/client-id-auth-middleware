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

func TestAuthorizationMiddlewareForAllCases(t *testing.T) {
	testCases := []struct {
		headers        map[string]string
		authMock       *mockClientAuthenticator
		expectedStatus int
		expectedCall   bool
		description    string
	}{
		{
			headers:        map[string]string{"Pass-Key": "some key"},
			authMock:       setupAuthMock("", "some key", errors.New("failed to authorize client")),
			expectedStatus: http.StatusUnauthorized,
			expectedCall:   false,
			description:    "Client ID Missing Unauthorized Scenario",
		},
		{
			headers:        map[string]string{"Client-ID": "client_id", "Pass-Key": "passkey"},
			authMock:       setupAuthMock("client_id", "passkey", nil),
			expectedStatus: http.StatusUnauthorized,
			expectedCall:   true,
			description:    "Authorization successful Scenario",
		},
		{
			headers:        map[string]string{"Client-ID": "client_id"},
			authMock:       setupAuthMock("client_id", "", errors.New("pass key is missing")),
			expectedStatus: http.StatusUnauthorized,
			expectedCall:   false,
			description:    "Pass Key Missing Scenario",
		},
		{
			headers:        map[string]string{"Client-ID": "some_client_id", "Pass-Key": "some_passkey"},
			authMock:       setupAuthMock("some_client_id", "some_passkey", errors.New("unauthorized request")),
			expectedStatus: http.StatusUnauthorized,
			expectedCall:   false,
			description:    "Client ID and PassKey available but unauthorized",
		},
	}

	var called bool
	nextHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	for _, tc := range testCases {
		r, w, err := setupRequest(http.MethodGet, "/authenticate", tc.headers)
		require.NoError(t, err, "should not have failed to create a request")

		authMw := NextAuthorizer(tc.authMock)

		authMw(w, r, nextHandlerFunc)

		assert.Equal(t, tc.expectedStatus, http.StatusUnauthorized, tc.description)
		assert.Equal(t, tc.expectedCall, called, tc.description)
		tc.authMock.AssertExpectations(t)
		called = false
	}
}

func setupRequest(method, url string, headers map[string]string) (*http.Request, *httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, nil, err
	}

	for hk, hval := range headers {
		r.Header.Add(hk, hval)
	}

	return r, httptest.NewRecorder(), nil
}

func setupAuthMock(clientID, passKey string, authErr error) *mockClientAuthenticator {
	mockClientAuthenticator := &mockClientAuthenticator{}
	mockClientAuthenticator.On("Authenticate", clientID, passKey).Return(authErr).Once()
	return mockClientAuthenticator
}
