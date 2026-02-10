package middleware

import (
	"net/http"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

func IPBlacklist(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil {
			c.Next()
			return
		}
		ip := c.ClientIP()
		blocked, _ := dnslog.IsIPBlacklistedWithContext(c.Request.Context(), ip)
		if blocked {
			response.Error(c, http.StatusForbidden, response.CodeForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
