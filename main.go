package main

import (
	"csv-reader/handlers"
	"csv-reader/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize services
	csvService := services.NewCSVService()

	// Initialize handlers
	handler := handlers.NewHandler(csvService)

	// Static file serving
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API routes
	http.HandleFunc("/", handler.HandleHome)
	http.HandleFunc("/api/headers", handler.HandleHeaders)
	http.HandleFunc("/api/query", handler.HandleQuery)
	http.HandleFunc("/api/upload", handler.HandleUpload)
	http.HandleFunc("/api/data", handler.HandleData)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
