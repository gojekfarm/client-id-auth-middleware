package clientauth

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type client struct {
	ClientID  string    `db:"client_id"`
	PassKey   string    `db:"pass_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type clientStore interface {
	GetClient(clientID string) (*client, error)
}

type clientRepository struct {
	db *sqlx.DB
}

func (r *clientRepository) GetClient(clientID string) (*client, error) {
	query := `
			SELECT client_id, pass_key
			FROM authorized_applications
			WHERE client_id = $1
			`
	authApplication := client{ClientID: clientID}

	tx := r.db
	err := tx.Get(&authApplication, query, clientID)
	if err != nil {
		return nil, err
	}

	return &authApplication, nil
}
