package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	learningbusiness "github.com/khoaphungnguyen/learning-tracker/internal/learning/business"
	learningstorage "github.com/khoaphungnguyen/learning-tracker/internal/learning/storage"
	learningtransport "github.com/khoaphungnguyen/learning-tracker/internal/learning/transport"
	middleware "github.com/khoaphungnguyen/learning-tracker/internal/middlewares"
	userbusiness "github.com/khoaphungnguyen/learning-tracker/internal/users/business"
	userstorage "github.com/khoaphungnguyen/learning-tracker/internal/users/storage"
	usertransport "github.com/khoaphungnguyen/learning-tracker/internal/users/transport"
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

	// Create a new user storage instance
	userDB, err := userstorage.NewUserStore()
	if err != nil {
		panic(err)
	}
	defer userDB.DB.Close()

	// Create a new learning storage instance
	learningDB, err := learningstorage.NewLearningStore()
	if err != nil {
		panic(err)
	}
	defer learningDB.DB.Close()

	// Create the tables if they don't exist
	err = userDB.CreateTable()
	if err != nil {
		panic(err)
	}
	// Create a new user service
	userService := userbusiness.NewUserService(userDB)
	userHandler := usertransport.NewUserHandler(userService, jwtKey)

	// Create a new learning service
	learningService := learningbusiness.NewLearningService(learningDB)
	learningHandler := learningtransport.NewLearningHandler(learningService)

	r := setupRouter(userHandler, learningHandler)
	r.Run(":8000")
}

func setupRouter(userHandler *usertransport.UserHandler, learningHandler *learningtransport.LearningHandler) *gin.Engine {
	r := gin.Default()
	// Create a new group for the API
	auth := r.Group("/auth")
	{
		// Add the login route
		auth.POST(("/login"), userHandler.Login)
		// Add the signup route
		auth.POST("/signup", userHandler.Signup)
		// Add the refresh token route
		auth.POST("/refresh", userHandler.RenewAccessToken)
	}

	// Create protected route
	protected := r.Group("/protected").Use(middleware.AuthMiddleware(userHandler.JWTKey))
	{
		// Get user profile
		protected.GET("/profile", userHandler.Profile)
		// Update user profile
		protected.PUT("/profile", userHandler.UpdateProfile)
		// Delete user profile
		protected.DELETE("/profile", userHandler.DeleteProfile)

		// Create a new goal
		protected.POST("/goals", learningHandler.CreateGoal)
		// Update a goal
		protected.PUT("/goals", learningHandler.UpdateGoal)
		// Delete a goal
		protected.DELETE("/goals/:id", learningHandler.DeleteGoal)
		// Get all goals
		protected.GET("/goals", learningHandler.GetAllGoalsByUserID)
		// Get a goal by ID
		protected.GET("/goals/:id", learningHandler.GetGoalByID)
		// Get all entries with the given goal ID
		protected.GET("/goals/:id/entries", learningHandler.GetAllEntriesByGoalID)

		// Create a new entry
		protected.POST("/entries", learningHandler.CreateEntry)
		// Update an entry
		protected.PUT("/entries", learningHandler.UpdateEntry)
		// Delete an entry
		protected.DELETE("/entries/:id", learningHandler.DeleteEntry)
		// Get an entry by ID
		protected.GET("/entries/:id", learningHandler.GetEntryByID)
		// Get all files by entry ID
		protected.GET("/entries/:id/files", learningHandler.GetAllFilesByEntryID)

		// Create a new file with the given entry ID
		protected.POST("/files", learningHandler.CreateFile)
		// Update a file
		protected.PUT("/files", learningHandler.UpdateFile)
		// Delete a file
		protected.DELETE("/files/:id", learningHandler.DeleteFile)
		// Get a file by ID
		protected.GET("/files/:id", learningHandler.GetFileByID)
		// Download a file
		protected.GET("/files/:id/download", learningHandler.DownloadFile)

	}

	return r
}
