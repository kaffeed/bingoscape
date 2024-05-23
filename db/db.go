package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	Db *sql.DB
}

func NewStore(dbName string) (Store, error) {
	db, err := getConnection(dbName)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db,
	}, nil
}

func getConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	// Init SQLite3 database
	db, err = sql.Open("pgx", dbName)
	if err != nil {
		// log.Fatalf("🔥 failed to connect to the database: %s", err.Error())
		return nil, fmt.Errorf("🔥 failed to connect to the database: %s", err)
	}

	log.Println("🚀 Connected Successfully to the Database")

	return db, nil
}
