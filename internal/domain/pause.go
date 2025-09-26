package domain

import (
	"dnslog_for_go/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StatusRequest 用于接收前端 JSON
type StatusRequest struct {
	Status string `json:"status"`
}

var paused bool // 全局状态

// InitPause 作为 Gin 的路由处理函数
func InitPause(c *gin.Context) {
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid pause/start request", zap.Error(err))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	PauseHandler(req) // 传入请求数据

	if paused {
		c.String(http.StatusOK, "System paused")
	} else {
		c.String(http.StatusOK, "System started")
	}
}

// PauseHandler 修改全局 paused 状态
func PauseHandler(req StatusRequest) {
	if req.Status == "pause" {
		paused = true
		log.Info("System paused")
	} else if req.Status == "start" {
		paused = false
		log.Info("System started")
	}
}

// IsPaused 对外提供
func IsPaused() bool {
	return paused
}
