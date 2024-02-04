package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestId string

		if c.GetHeader("X-Request-Id") != "" {
			requestId = c.GetHeader("X-Request-Id")
		} else {
			requestId = uuid.NewString()
		}

		c.Request.Header.Set("X-Request-Id", requestId)

		c.Header("X-Request-Id", requestId)

		c.Next()
	}
}
