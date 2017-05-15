package clientauth

import "github.com/stretchr/testify/mock"

type mockClientAuthenticator struct {
	mock.Mock
}

func (mca *mockClientAuthenticator) Authenticate(clientID, passKey string) error {
	args := mca.Called(clientID, passKey)

	return args.Error(0)
}

type mockClientRepository struct {
	mock.Mock
}

func (mcr *mockClientRepository) getClient(clientID string) (*client, error) {
	args := mcr.Called(clientID)
	var myClient *client

	if args[0] != nil {
		myClient = args[0].(*client)
	}

	return myClient, args.Error(1)
}
