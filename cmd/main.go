package main

import (
	"log"
	"os"

	"finance-backend/internal/database"
	"finance-backend/internal/handlers"
	"finance-backend/internal/repository"
	"finance-backend/internal/middleware"
    "fmt"
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

	// Start server
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
    userRepo := &repository.UserRepository{DB: database.DB}
    authHandler := &handlers.AuthHandler{Repo: userRepo}

    r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
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
	r.Run(":" + port)
	
}