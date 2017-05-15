package clientauth

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatesAndFindsAnAuthorizedApplication(t *testing.T) {
	dbDriver := "postgres"
	dbConnURL := "dbname=client_auth user=postgres host=localhost sslmode=disable"

	db := loadDatabase(dbDriver, dbConnURL)
	db.Exec("INSERT INTO authorized_applications (client_id, pass_key) VALUES ('DUMMY-CLIENT-ID', 'DUMMY-PASSKEY')")

	authorizedApplication := client{
		ClientID: "DUMMY-CLIENT-ID",
		PassKey:  "DUMMY-PASSKEY",
	}

	repository := &clientRepository{db: db}
	found, err := repository.getClient(authorizedApplication.ClientID)
	require.NoError(t, err, "should not have failed to get client from db")

	assert.Equal(t, authorizedApplication.ClientID, found.ClientID)
	assert.Equal(t, authorizedApplication.PassKey, found.PassKey)

	db.Exec("DELETE FROM authorized_applications WHERE client_id = 'DUMMY-CLIENT-ID' ")
}

func TestGetClientFailsWhenNoRecords(t *testing.T) {
	dbDriver := "postgres"
	dbConnURL := "dbname=client_auth user=postgres host=localhost sslmode=disable"

	db := loadDatabase(dbDriver, dbConnURL)

	authorizedApplication := client{
		ClientID: "DUMMY-CLIENT-ID",
		PassKey:  "DUMMY-PASSKEY",
	}

	repository := &clientRepository{db: db}
	found, err := repository.getClient(authorizedApplication.ClientID)
	require.Error(t, err, "should have failed to get client from db")

	assert.Nil(t, found)
}
