package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/khoaphungnguyen/learning-tracker/internal/business"
	"github.com/khoaphungnguyen/learning-tracker/internal/storage"
	"github.com/khoaphungnguyen/learning-tracker/internal/transport"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch JWT key from the environment
	jwtKey := os.Getenv("JWT_SECRET_KEY")
	if jwtKey == "" {
		log.Fatal("JWT_SECRET_KEY not set in .env file")
	}

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
	learningHandler := transport.NewNetHandler(learningService, jwtKey)

	// Initialzie the router
	router := http.NewServeMux()

	// Define API routes using handlers from the "api" package
	learningHandler.SetupRoutes(router)

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
