package dnslog

import (
	"github.com/genwilliam/dnslog_for_go/pkg/response"
	"github.com/genwilliam/dnslog_for_go/pkg/utils"

	"strconv"

	"github.com/gin-gonic/gin"
)

// ListRecordsHandler 返回当前捕获到的 DNS 请求记录（支持分页）
func ListRecordsHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	records := GetRecords()
	paged := utils.Paginate(records, page, pageSize)

	response.Success(c, gin.H{
		"items": paged,
		"total": len(records),
		"page":  page,
	})
}
