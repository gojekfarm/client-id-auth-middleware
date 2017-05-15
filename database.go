package clientauth

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func loadDatabase(dbDriver, dbConnURL string) *sqlx.DB {
	return sqlx.MustConnect(dbDriver, dbConnURL)
}
