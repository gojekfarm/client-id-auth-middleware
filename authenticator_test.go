package clientauth

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientAuthenticationSuite struct {
	suite.Suite
}

func (suite *ClientAuthentication) SetupSuite() {
}

func TestClientAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(ClientAuthenticationSuite))
}
