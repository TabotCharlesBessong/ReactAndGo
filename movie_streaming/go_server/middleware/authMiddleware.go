package middleware

import (
	"net/http"

	"github.com/TabotCharlesBessong/ReactAndGo/tree/movie_streamer/movie_streaming/go_server/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// c.Set("email", claims.Email)
		// c.Set("firstName", claims.FirstName)
		// c.Set("lastName", claims.LastName)
		c.Set("role", claims.Role)
		c.Set("userId", claims.UserId)

		c.Next()
	}
}
