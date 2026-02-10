package dnslog

import (
	"strconv"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListRecordsHandler 返回当前捕获到的 DNS 请求记录（支持分页）
func ListRecordsHandler(c *gin.Context) {
	cfg := config.Get()
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", strconv.Itoa(cfg.DefaultPageSize))
	order := c.DefaultQuery("order", "desc")
	cursorStr := c.DefaultQuery("cursor", "")

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
		Domain:   c.Query("domain"),
		ClientIP: c.Query("client_ip"),
		Protocol: c.Query("protocol"),
		QType:    c.Query("qtype"),
		Token:    c.Query("token"),
		Order:    order,
	}
	if cursorStr != "" {
		if v, err := strconv.ParseInt(cursorStr, 10, 64); err == nil {
			filter.Cursor = v
		}
	}

	if start := c.Query("start"); start != "" {
		if v, err := strconv.ParseInt(start, 10, 64); err == nil {
			filter.Start = v
		}
	}
	if end := c.Query("end"); end != "" {
		if v, err := strconv.ParseInt(end, 10, 64); err == nil {
			filter.End = v
		}
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
	})
}
