package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/ramzeng/ai-endpoint/package/error"
)

func Timeout() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(60*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": error.GatewayTimeout("request timeout"),
			})
		}),
	)
}
