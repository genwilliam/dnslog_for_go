package dnslog

import (
	"net/http"
	"time"

	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

// SetTokenWebhookHandler 绑定 token webhook（仅支持 FIRST_HIT）
func SetTokenWebhookHandler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	var req struct {
		URL    string `json:"webhook_url" binding:"required"`
		Secret string `json:"secret"`
		Mode   string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	if req.Mode == "" {
		req.Mode = "FIRST_HIT"
	}
	if req.Mode != "FIRST_HIT" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	if err := UpsertTokenWebhookWithContext(c.Request.Context(), token, req.URL, req.Secret, req.Mode, time.Now().UnixMilli()); err != nil {
		if err == ErrSecretKeyRequired {
			response.Error(c, http.StatusBadRequest, response.CodeWebhookSecretKeyRequired)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"token": token, "webhook_url": req.URL, "mode": req.Mode})
}

// GetTokenWebhookHandler 获取 token webhook
func GetTokenWebhookHandler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	hook, err := GetTokenWebhookWithContext(c.Request.Context(), token)
	if err == ErrWebhookNotFound {
		response.Error(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{
		"token":       hook.Token,
		"webhook_url": hook.URL,
		"mode":        hook.Mode,
		"enabled":     hook.Enabled,
		"created_at":  hook.CreatedAt,
	})
}

// DisableTokenWebhookHandler 禁用 token webhook
func DisableTokenWebhookHandler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	if err := DisableTokenWebhookWithContext(c.Request.Context(), token); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"token": token, "disabled": true})
}
