package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TraceID   string      `json:"trace_id"`
	Timestamp int64       `json:"timestamp"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	traceID := GetTraceID(c)

	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   "ok",
		Data:      data,
		TraceID:   traceID,
		Timestamp: time.Now().UnixMilli(),
	})
}

// Error 返回错误响应（不带 trace）
func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:      httpCode,
		Message:   msg,
		TraceID:   "",
		Timestamp: time.Now().UnixMilli(),
	})
}

// ErrorWithTrace 返回带 trace_id 的错误响应
func ErrorWithTrace(c *gin.Context, httpCode int, msg string, traceID string) {
	c.JSON(httpCode, Response{
		Code:      httpCode,
		Message:   msg,
		TraceID:   traceID,
		Timestamp: time.Now().UnixMilli(),
	})
}

// GetTraceID 从 context 获取 trace_id（如果有中间件）
func GetTraceID(c *gin.Context) string {
	v, exists := c.Get("trace_id")
	if !exists {
		return ""
	}
	traceID, ok := v.(string)
	if !ok {
		return ""
	}
	return traceID
}
