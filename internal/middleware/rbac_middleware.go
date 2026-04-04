package middleware
import (
	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(403, gin.H{"error": "No role found"})
			c.Abort()
			return
		}

		userRole := roleVal.(string)

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
		c.Abort()
	}
}