package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	db, err := userstorage.NewUserStore()
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Create the tables if they don't exist
	err = db.CreateTable()
	if err != nil {
		panic(err)
	}
	// Create a new user service
	userService := userbusiness.NewUserService(db)
	userHandler := usertransport.NewUserHandler(userService, jwtKey)

	r := setupRouter(userHandler)
	r.Run(":8000")
}

func setupRouter(userHandler *usertransport.UserHandler) *gin.Engine {
	r := gin.Default()
	// Create a new group for the API
	auth := r.Group("/auth")
	{
		// Add the login route
		auth.POST(("/login"), userHandler.Login)
		// Add the signup route
		auth.POST("/signup", userHandler.Signup)
	}

	// Create protected route
	protected := r.Group("/protected").Use(middleware.AuthMiddleware(userHandler.JWTKey))
	{
		protected.GET("/profile", userHandler.Profile)
	}
	// Return the router
	return r
}
