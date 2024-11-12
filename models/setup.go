package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

//var DB *gorm.DB

//func ConnectDatabase() {
//	database, err := gorm.Open(sqlite.Open("storage/usagestats.db"), &gorm.Config{})
//
//	if err != nil {
//		panic("Failed to connect to database!")
//	}
//
//	DB = database
//}
//
//func DBMigrate() {
//	DB.AutoMigrate(&UserStats{})
//}

func init() {
	db, err := sql.Open("sqlite3", "storage/usagestats.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
