package middleware

import (
	"net/http"
	"strings"

	utils "golang.elasticsearch/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header value
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized Access."})
			c.Abort() // Stop further processing
			return
		}

		// Check if the header format is "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized Access."})
			c.Abort() // Stop further processing
			return
		}

		// Extract the token string by trimming the "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Bearer Token."})
			c.Abort()
			return
		}

		// store the token or relevant user info in the context for handlers
		c.Set("authToken", token)

		// Continue to the next handler
		c.Next()
	}
}
