package middleware

import (
	"github.com/genwilliam/dnslog_for_go/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TraceID sets a trace_id for the request lifecycle.
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := utils.GenerateTraceID()
		c.Set("trace_id", traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
