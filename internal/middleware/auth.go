package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ramzeng/ai-endpoint/internal/client"
	"github.com/ramzeng/ai-endpoint/package/database"
	"github.com/ramzeng/ai-endpoint/package/error"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": error.Unauthorized(),
			})
			return
		}

		secret := strings.TrimPrefix(token, "Bearer ")

		var requestClient client.Client

		results := database.Where("secret = ?", secret).First(&requestClient)

		if results.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": error.Unauthorized(),
			})
			return
		}

		c.Set("client", requestClient)
		c.Next()
	}
}
