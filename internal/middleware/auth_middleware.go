package middleware

import (
	"net/http"
	"strings"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")

		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
        
		parts := strings.SplitN(authHeader, " ", 2)

if len(parts) != 2 || parts[0] != "Bearer" {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Invalid token format",
	})
	c.Abort()
	return
}

        tokenString := parts[1]
		

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		c.Next()

	}
}