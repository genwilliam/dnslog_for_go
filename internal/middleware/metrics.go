package middleware

import (
	"strconv"

	"github.com/genwilliam/dnslog_for_go/internal/metrics"

	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		metrics.APIRequestsTotal.WithLabelValues(path, c.Request.Method, strconv.Itoa(status)).Inc()
	}
}
