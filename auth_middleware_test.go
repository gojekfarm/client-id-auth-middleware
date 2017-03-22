package clientIdAuth_test

import (
	"client-id-auth-middleware"
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

func setUp() {
}

type MockClientRepositoryErrors struct{}

type MockClientRepository struct{}

func (m *MockClientRepository) GetClient(clientId string) (*clientIdAuth.Client, error) {
	return &clientIdAuth.Client{
		ClientId: "ClientID",
		PassKey:  "Pass-Key",
	}, nil
}

func (m *MockClientRepositoryErrors) GetClient(clientId string) (*clientIdAuth.Client, error) {
	return nil, errors.New("no row in db for this client id")
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

func TestWithClientIdAndPassKeyAuthorizationReturnUnauthorizedIfClientIdIsMissing(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepository{}

	ts := httptest.NewServer(clientIdAuth.WithClientIdAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Pass-Key", "some key")
	response, _ := client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestWithClientIdAndPassKeyAuthorizationReturnUnauthorizedIfPassKeyIsMissing(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepository{}
	ts := httptest.NewServer(clientIdAuth.WithClientIdAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})))
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "some client id")
	response, _ := client.Do(req)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestWithClientIdAndPassKeyAuthorizationReturnUnauthorizedIsWrongClientIdAndPassKeyIsSent(t *testing.T) {
	setUp()
	db := initDB()
	db.Exec("INSERT INTO authorised_applications (client_id, pass_key) VALUES ('DUMMY-CLIENT-ID', 'DUMMY-PASSKEY')")

	clientRepository := &MockClientRepository{}
	ts := httptest.NewServer(clientIdAuth.WithClientIdAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "DUMMY-CLIENT-ID")
	req.Header.Add("Pass-Key", "some pass key")

	response, _ := client.Do(req)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	db.Exec("DELETE FROM authorised_applications WHERE client_id = 'DUMMY-CLIENT-ID' ")
}

func TestWithClientIdAndPassKeyAuthorizationReturnUnAuthorizedIfClientIdIsNotPresentInDB(t *testing.T) {
	setUp()
	clientRepository := &MockClientRepositoryErrors{}
	ts := httptest.NewServer(clientIdAuth.WithClientIdAndPassKeyAuthorization(clientRepository)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)

	client := &http.Client{}
	req.Header.Add("Client-ID", "DUMMY-CLIENT-ID")
	req.Header.Add("Pass-Key", "some pass key")

	response, _ := client.Do(req)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
