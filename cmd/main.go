package main

import (
	"fmt"
	"time"

	"github.com/khoaphungnguyen/learning-tracker/internal/database"
)

func main() {

	dbService, err := database.NewDBService()
	if err != nil {
		panic(err)
	}
	defer dbService.DB.Close()

	// Create the tables
	err = dbService.CreateTable()
	if err != nil {
		panic(err)
	}

	// Create new learning goals
	err = dbService.CreateGoal("Learning Algorithms", time.Now(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		panic(err)
	}

	// Get a goal by ID
	goal, err := dbService.GetGoal(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Goal By ID:", goal)

	// Get all learning goals
	goals, err := dbService.GetAllGoals()
	if err != nil {
		panic(err)
	}
	fmt.Println("All Goals:", goals)

	// Create new learning entries by goal ID
	err = dbService.CreateEntry(1, "Learning Binary Search", "O(logn) time", time.Now(), false)
	if err != nil {
		panic(err)
	}

	// Get a learning entry by ID
	entry, err := dbService.GetEntry(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Entry By ID:", entry)

	// Create new learning entries

	// // Initialzie the router
	// router := http.NewServeMux()

	// // Define API routes using handlers from the "api" package
	// api.SetupRoutes(router)

	// // Start the server
	// port := ":8000"

	// server := &http.Server{
	// 	Addr:    port,
	// 	Handler: router,
	// }
	// fmt.Printf("Server listening on port %s\n", port)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	panic(err)
	// }

}
