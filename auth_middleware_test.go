package clientauth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/irfn/goconfig"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type MockClientRepositoryErrors struct{}

func (m *MockClientRepositoryErrors) GetClient(clientID string) (*client, error) {
	return nil, errors.New("no row in db for this client id")
}

type MockClientRepository struct{}

func (m *MockClientRepository) GetClient(clientID string) (*client, error) {
	return &client{
		ClientID: "ClientID",
		PassKey:  "Pass-Key",
	}, nil
}

func initDB() *sqlx.DB {
	dbConf := goconfig.LoadDbConf()
	db, err := sqlx.Connect(dbConf.Driver(), dbConf.Url())
	if err != nil {
		log.Panic(fmt.Errorf("Unable to connect to the DB: %v", err))
	}

	db.SetMaxOpenConns(dbConf.MaxConn())
	db.SetMaxIdleConns(dbConf.IdleConn())
	db.SetConnMaxLifetime(dbConf.ConnMaxLifetime())
	return db
}

func TestWithClientIDAndPassKeyAuthorizationReturnUnauthorizedIfClientIDIsMissing(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepository{}

	ts := httptest.NewServer(WithClientIDAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Pass-Key", "some key")
	response, _ := client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestWithClientIDAndPassKeyAuthorizationReturnUnauthorizedIfPassKeyIsMissing(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepository{}
	ts := httptest.NewServer(WithClientIDAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "some client id")
	response, _ := client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestWithClientIDAndPassKeyAuthorizationReturnUnauthorizedIsWrongClientIDAndPassKeyIsSent(t *testing.T) {
	setUp()
	db := initDB()
	db.Exec("INSERT INTO authorized_applications (client_id, pass_key) VALUES ('DUMMY-CLIENT-ID', 'DUMMY-PASSKEY')")

	clientRepository := &MockClientRepository{}
	ts := httptest.NewServer(WithClientIDAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "DUMMY-CLIENT-ID")
	req.Header.Add("Pass-Key", "some pass key")

	response, _ := client.Do(req)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	db.Exec("DELETE FROM authorized_applications WHERE client_id = 'DUMMY-CLIENT-ID' ")
}

func TestWithClientIDAndPassKeyAuthorizationReturnUnAuthorizedIfClientIDIsNotPresentInDB(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepositoryErrors{}
	ts := httptest.NewServer(WithClientIDAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "DUMMY-CLIENT-ID")
	req.Header.Add("Pass-Key", "some pass key")

	response, _ := client.Do(req)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
