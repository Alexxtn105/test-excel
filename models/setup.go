package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "storage/usagestats.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
