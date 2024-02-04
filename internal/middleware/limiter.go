package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ramzeng/ai-endpoint/internal/client"
	"github.com/ramzeng/ai-endpoint/package/error"
	"github.com/ramzeng/ai-endpoint/package/logger"
	"github.com/ramzeng/ai-endpoint/package/redis"
	"github.com/ramzeng/ai-endpoint/package/toolkit"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type limiterRequest struct {
	Model   string
	Version string
}

func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodHead || c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		client := c.MustGet("client").(client.Client)

		body, _ := toolkit.ReadRequestBody(c.Request)

		var request limiterRequest

		_ = binding.JSON.BindBody(body, &request)

		if request.Model == "" && request.Version == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": error.BadRequest("resource parameter missing"),
			})
			return
		}

		for _, rateLimit := range client.RateLimits {
			if rateLimit.Model != request.Model && rateLimit.Model != request.Version {
				continue
			}

			if rateLimit.MaxRequestsPerMinute <= 0 {
				continue
			}

			key := fmt.Sprintf(
				"limiter_%d_%s_%d",
				client.Id,
				rateLimit.Model,
				rateLimit.MaxRequestsPerMinute,
			)

			hits, err := redis.IncrNX(key, time.Minute)

			if err != nil {
				logger.Error(
					"app",
					"[Limiter]: hits increment failed",
					zap.String("event", "request_limiter_error"),
					zap.Error(err),
				)
				c.Next()
				return
			}

			ttl, err := redis.TTL(key)

			if err != nil {
				logger.Error(
					"app",
					"[Limiter]: get hits ttl failed",
					zap.String("event", "request_limiter_error"),
					zap.Error(err),
				)
				c.Next()
				return
			}

			var remaining uint64

			if rateLimit.MaxRequestsPerMinute > hits {
				remaining = rateLimit.MaxRequestsPerMinute - hits
			}

			c.Header("X-RateLimit-Limit", cast.ToString(rateLimit.MaxRequestsPerMinute))
			c.Header("X-RateLimit-Remaining", cast.ToString(remaining))
			c.Header("X-RateLimit-Reset", cast.ToString(time.Now().Add(ttl).Unix()))

			if hits <= rateLimit.MaxRequestsPerMinute {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": error.TooManyRequests(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": error.Forbidden("no resource access permission"),
		})
	}
}
