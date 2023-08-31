package main

import (
	"fmt"
	"net/http"

	"github.com/khoaphungnguyen/learning-tracker/internal/business"
	"github.com/khoaphungnguyen/learning-tracker/internal/storage"
	"github.com/khoaphungnguyen/learning-tracker/internal/transport"
)

func main() {
	// Create a new storage instance
	db, err := storage.NewLearningStore()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()
	// Create the tables if they don't exist
	err = db.CreateTable()
	if err != nil {
		panic(err)
	}
	// Create a new learning service
	learningService := business.NewLearningService(db)

	// Create a new transport hanlder
	handler := transport.NewNetHandler(learningService)

	// Initialzie the router
	router := http.NewServeMux()

	// Define API routes using handlers from the "api" package
	handler.SetupRoutes(router)

	// Start the HTTP server
	port := ":8000"
	server := &http.Server{
		Addr:    port,
		Handler: router,
	}
	fmt.Printf("Server listening on port %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
