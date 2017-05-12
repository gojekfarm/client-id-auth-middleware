package clientauth

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func loadDatabase(dbConnURL, dbDriver string) *sqlx.DB {
	db, err := sqlx.Connect(dbDriver, dbConnURL)
	if err != nil {
		log.Panic(fmt.Errorf("Unable to connect to the DB: %v", err))
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(20)

	return db
}
