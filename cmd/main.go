package main

import (
	"fmt"
	"net/http"

	"github.com/khoaphungnguyen/learning-tracker/internal/api"
)

func main() {
	// Initialzie the router
	router := http.NewServeMux()

	// Define API routes using handlers from the "api" package
	api.SetupRoutes(router)

	// Start the server
	port := ":8000"

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}
	fmt.Printf("Server listening on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
