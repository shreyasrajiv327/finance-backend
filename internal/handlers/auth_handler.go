package handlers

import (
	"finance-backend/internal/models"
	"finance-backend/internal/repository"
	"net/http"
	"os"
    "time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)


type AuthHandler struct{
	Repo *repository.UserRepository
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput

	// Parse request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	// Validation
	if input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password required",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Convert to User model
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "viewer", // default role
	}

	// Save to DB
	err = h.Repo.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create user (email might exist)",
		})
		return
	}

	// Return response (safe)
	c.JSON(http.StatusCreated, user)
}


func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.Repo.GetByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 🔥 FIX: CHECK bcrypt result
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
     fmt.Println("INPUT EMAIL:", input.Email)
fmt.Println("INPUT PASS:", input.Password)
fmt.Println("DB HASH:", user.Password)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}