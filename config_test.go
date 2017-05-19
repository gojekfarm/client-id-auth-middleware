package clientauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigHeaderOverriding(t *testing.T) {
	config := NewConfig("db_driver", "conn_url")

	config.SetHeaderConfig("client_id", "pass_key")

	assert.Equal(t, "client_id", config.HeaderConfig.ClientIDName)
	assert.Equal(t, "pass_key", config.HeaderConfig.PassKeyName)
}

func TestDefaultConfigHeader(t *testing.T) {
	config := NewConfig("db_driver", "conn_url")

	assert.Equal(t, "Client-ID", config.HeaderConfig.ClientIDName)
	assert.Equal(t, "Pass-Key", config.HeaderConfig.PassKeyName)
}
