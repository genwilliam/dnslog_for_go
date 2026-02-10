package dnslog

import (
	"net/http"
	"strconv"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetTokenStatusHandler 返回 token 状态（主入口）
func GetTokenStatusHandler(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	ts, err := GetTokenStatusWithContext(c.Request.Context(), token)
	if err == ErrTokenNotFound {
		response.Error(c, http.StatusNotFound, response.CodeTokenNotFound)
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}

	nowMs := time.Now().UnixMilli()
	if ts.Status != "EXPIRED" && ts.ExpiresAt > 0 && nowMs > ts.ExpiresAt {
		if _, err := MaybeExpireTokenWithContext(c.Request.Context(), token, nowMs); err == nil {
			ts.Status = "EXPIRED"
			ts.UpdatedAt = nowMs
		}
	}

	response.Success(c, gin.H{
		"token":      ts.Token,
		"domain":     ts.Domain,
		"status":     ts.Status,
		"first_seen": ts.FirstSeen,
		"last_seen":  ts.LastSeen,
		"hit_count":  ts.HitCount,
		"expires_at": ts.ExpiresAt,
		"expired":    ts.Status == "EXPIRED",
	})
}

// GetTokenRecordsHandler 返回 token 的原始记录
func GetTokenRecordsHandler(c *gin.Context) {
	cfg := config.Get()
	token := c.Param("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", strconv.Itoa(cfg.DefaultPageSize))
	order := c.DefaultQuery("order", "desc")

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

	filter := ListFilter{
		Page:     page,
		PageSize: pageSize,
		Token:    token,
		Order:    order,
	}

	items, total, err := ListRecordsWithContext(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  pageSize,
		"order": filter.Order,
	})
}

// ListTokensHandler 返回 token 列表
func ListTokensHandler(c *gin.Context) {
	cfg := config.Get()
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", strconv.Itoa(cfg.DefaultPageSize))
	order := c.DefaultQuery("order", "desc")
	orderBy := c.DefaultQuery("orderBy", "created_at")

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

	filter := TokenListFilter{
		Page:        page,
		PageSize:    pageSize,
		Status:      c.Query("status"),
		Order:       order,
		OrderBy:     orderBy,
	}
	if v := c.Query("keyword"); v != "" {
		filter.Keyword = v
	}

	if v := c.Query("created_start"); v != "" {
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.CreatedStart = ts
		}
	}
	if v := c.Query("created_end"); v != "" {
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.CreatedEnd = ts
		}
	}
	if v := c.Query("last_start"); v != "" {
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.LastStart = ts
		}
	}
	if v := c.Query("last_end"); v != "" {
		if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.LastEnd = ts
		}
	}

	items, total, err := ListTokensWithContext(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, 500, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  pageSize,
		"order": filter.Order,
	})
}
