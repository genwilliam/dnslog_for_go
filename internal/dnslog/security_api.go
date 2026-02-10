package dnslog

import (
	"net/http"
	"strconv"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateAPIKeyHandler 创建 API Key（返回明文 key，仅一次）
func CreateAPIKeyHandler(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	plain, hash, err := GenerateAPIKey()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	cfg := config.Get()
	nowMs := time.Now().UnixMilli()
	if cfg != nil && cfg.APIKeyRequired {
		if _, ok := c.Get("api_key_id"); !ok {
			id, err := CreateBootstrapAPIKeyWithContext(c.Request.Context(), req.Name, hash, req.Comment, nowMs)
			if err == ErrBootstrapConflict {
				response.Error(c, http.StatusConflict, response.CodeAPIKeyAlreadyInitialized)
				return
			}
			if err != nil {
				response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
				return
			}
			response.Success(c, gin.H{
				"id":   id,
				"name": req.Name,
				"key":  plain,
			})
			return
		}
	}

	id, err := CreateAPIKeyWithContext(c.Request.Context(), req.Name, hash, req.Comment, nowMs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"id":   id,
		"name": req.Name,
		"key":  plain,
	})
}

// ListAPIKeysHandler 列出 API Keys
func ListAPIKeysHandler(c *gin.Context) {
	cfg := config.Get()
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", strconv.Itoa(cfg.DefaultPageSize))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = cfg.DefaultPageSize
	}
	if pageSize > cfg.MaxPageSize {
		pageSize = cfg.MaxPageSize
	}

	items, total, err := ListAPIKeysWithContext(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	respItems := make([]gin.H, 0, len(items))
	for _, k := range items {
		hashPrefix := ""
		if len(k.APIKey) >= 6 {
			hashPrefix = k.APIKey[:6]
		}
		respItems = append(respItems, gin.H{
			"id":           k.ID,
			"name":         k.Name,
			"enabled":      k.Enabled,
			"created_at":   k.CreatedAt,
			"last_used_at": k.LastUsedAt,
			"comment":      k.Comment,
			"hash_prefix":  hashPrefix,
		})
	}

	response.Success(c, gin.H{
		"items": respItems,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// DisableAPIKeyHandler 禁用 API Key
func DisableAPIKeyHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	if err := SetAPIKeyEnabledWithContext(c.Request.Context(), id, false); err != nil {
		if err == ErrAPIKeyNotFound {
			response.Error(c, http.StatusNotFound, response.CodeNotFound)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"id": id, "disabled": true})
}

// CreateAPIKeyWithBootstrapHandler 紧急恢复：需要 bootstrap token
func CreateAPIKeyWithBootstrapHandler(c *gin.Context) {
	cfg := config.Get()
	if cfg == nil || !cfg.BootstrapEnabled {
		response.Error(c, http.StatusForbidden, response.CodeForbidden)
		return
	}
	token := c.GetHeader("X-Bootstrap-Token")
	if token == "" || cfg.BootstrapToken == "" || token != cfg.BootstrapToken {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized)
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	plain, hash, err := GenerateAPIKey()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	nowMs := time.Now().UnixMilli()
	id, err := CreateAPIKeyWithContext(c.Request.Context(), req.Name, hash, req.Comment, nowMs)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"id":   id,
		"name": req.Name,
		"key":  plain,
	})
}

// AddBlacklistHandler 添加 IP 黑名单
func AddBlacklistHandler(c *gin.Context) {
	var req struct {
		IP     string `json:"ip" binding:"required"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	if err := AddBlacklistIPWithContext(c.Request.Context(), req.IP, req.Reason, time.Now().UnixMilli()); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"ip": req.IP, "enabled": true})
}

// ListBlacklistHandler 列出黑名单
func ListBlacklistHandler(c *gin.Context) {
	cfg := config.Get()
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", strconv.Itoa(cfg.DefaultPageSize))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = cfg.DefaultPageSize
	}
	if pageSize > cfg.MaxPageSize {
		pageSize = cfg.MaxPageSize
	}

	items, total, err := ListBlacklistWithContext(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// DisableBlacklistHandler 禁用黑名单条目
func DisableBlacklistHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}
	if err := DisableBlacklistIPWithContext(c.Request.Context(), id); err != nil {
		if err == ErrBlacklistNotFound {
			response.Error(c, http.StatusNotFound, response.CodeNotFound)
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"id": id, "disabled": true})
}
