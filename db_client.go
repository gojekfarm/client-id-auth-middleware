package clientauth

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	ClientID  string    `db:"client_id"`
	PassKey   string    `db:"pass_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type clientStore interface {
	GetClient(clientID string) (*Client, error)
}

type ClientRepository struct {
	DB *sqlx.DB
}

func (r *ClientRepository) GetClient(clientID string) (*Client, error) {
	query := `
			SELECT client_id, pass_key
			FROM authorized_applications
			WHERE client_id = $1
			`
	authApplication := Client{ClientID: clientID}

	tx := r.DB

	err := tx.Get(&authApplication, query, clientID)
	if err != nil {
		return nil, err
	}
	return &authApplication, nil
}
