package clientauth

import "github.com/stretchr/testify/mock"

type MockClientAuthenticator struct {
	mock.Mock
}

func (mca *MockClientAuthenticator) Authenticate(clientID, passKey string) error {
	args := mca.Called(clientID, passKey)

	return args.Error(0)
}
