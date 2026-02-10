package middleware

import (
	"database/sql"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Audit(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil || !cfg.AuditEnabled {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		traceID := utils.GenerateTraceID()
		if v, ok := c.Get("trace_id"); ok {
			if t, ok := v.(string); ok && t != "" {
				traceID = t
			}
		}

		var apiKeyID sql.NullInt64
		if v, ok := c.Get("api_key_id"); ok {
			if id, ok := v.(int64); ok {
				apiKeyID = sql.NullInt64{Int64: id, Valid: true}
			}
		}

		token := c.Param("token")
		if token == "" {
			token = c.Query("token")
		}

		_ = dnslog.EnqueueAuditLog(dnslog.AuditLog{
			TraceID:    traceID,
			APIKeyID:   apiKeyID,
			Path:       c.FullPath(),
			Method:     c.Request.Method,
			ClientIP:   c.ClientIP(),
			StatusCode: status,
			LatencyMs:  latency,
			Token:      token,
			CreatedAt:  time.Now().UnixMilli(),
		})
	}
}
