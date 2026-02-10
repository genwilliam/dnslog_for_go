package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/infra"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

func RateLimit(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil || !cfg.RateLimitEnabled {
			c.Next()
			return
		}

		client := infra.GetRedis()
		if client == nil {
			response.Error(c, http.StatusServiceUnavailable, response.CodeRateLimitOff)
			c.Abort()
			return
		}

		scope := c.GetHeader(apiKeyHeader)
		if scope == "" {
			scope = c.ClientIP()
		}

		key := fmt.Sprintf("rl:%s:%s", c.FullPath(), scope)
		limit := cfg.RateLimitMaxRequests
		window := time.Duration(cfg.RateLimitWindowSeconds) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			response.Error(c, http.StatusServiceUnavailable, response.CodeRateLimitError)
			c.Abort()
			return
		}
		if count == 1 {
			_ = client.Expire(ctx, key, window).Err()
		}
		if int(count) > limit {
			response.Error(c, http.StatusTooManyRequests, response.CodeRateLimited)
			c.Abort()
			return
		}

		c.Next()
	}
}
