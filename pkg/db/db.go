package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(filepath string) error {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
