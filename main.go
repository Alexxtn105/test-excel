package main

import (
	"log"
	"net/http"
	"test-excel/controllers"
)

// addr The server port
var addr = ":8088"

func main() {
	//fmt.Println("hello")

	// Создаем пустую БД
	//models.ConnectDatabase()
	//models.DBMigrate()

	// Register the handler function for the "/hello" endpoint
	http.HandleFunc("GET /report", controllers.ReportHandler)
	http.HandleFunc("GET /report/csv", controllers.ReportCSVHandler)
	http.HandleFunc("GET /upload/csv", controllers.ReportCSVDownloadHandler)

	log.Printf("Starting server on %s\n", addr)

	// Start the server
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
