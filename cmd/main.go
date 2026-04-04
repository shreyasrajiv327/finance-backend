package main

import (
	"log"
	"os"
	"fmt"

	"finance-backend/internal/database"
	"finance-backend/internal/handlers"
	"finance-backend/internal/repository"
	"finance-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Println("No .env file found")
		}
	}

	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))

	// Connect DB
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	database.Migrate()

	r := gin.Default()

	// Public route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Setup repos + handlers
	userRepo := &repository.UserRepository{DB: database.DB}
	authHandler := &handlers.AuthHandler{Repo: userRepo}

	recordRepo := &repository.RecordRepository{DB: database.DB}
	recordHandler := &handlers.RecordHandler{Repo: recordRepo}

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Protected routes (ONLY ONE GROUP)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")

		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   email,
		})
	})

	protected.POST("/records", recordHandler.CreateRecord)

	// Start server (ALWAYS LAST)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}