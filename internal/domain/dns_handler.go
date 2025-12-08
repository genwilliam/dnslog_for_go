package domain

import (
	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/domain/dns_server"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"github.com/genwilliam/dnslog_for_go/pkg/response"
	"github.com/genwilliam/dnslog_for_go/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ChangeDNSRequest 修改DNS请求体
type ChangeDNSRequest struct {
	Num int `json:"num"`
}

// ChangePactRequest 修改协议请求体
type ChangePactRequest struct {
	Pact string `json:"pact"`
}

// DNSQueryResult DNS 查询结果结构体
type DNSQueryResult struct {
	Domain  string   `json:"domain"`
	Results []string `json:"results"` // 存储多个结果
}

// ShowForm 展示表单
func ShowForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// SubmitDomain 提交域名并查询
func SubmitDomain(c *gin.Context) {

	traceID := utils.GenerateTraceID()
	clientIP := c.ClientIP()

	if IsPaused() {
		log.Warn("系统暂停，拒绝查询",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
		)
		response.ErrorWithTrace(c, 503, "系统已暂停，无法查询域名", traceID)
		return
	}

	var req struct {
		DomainName string `json:"domain_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("参数格式错误",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.Error(err),
		)
		response.ErrorWithTrace(c, 400, "参数格式错误: "+err.Error(), traceID)
		return
	}

	if req.DomainName == "" {
		response.ErrorWithTrace(c, 400, "域名不能为空", traceID)
		return
	}

	if !utils.StandardizeDomain(req.DomainName) {
		log.Info("域名格式不合法",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.String("domain", req.DomainName),
		)
		response.ErrorWithTrace(c, 400, "域名不合法，请重新输入", traceID)
		return
	}

	// 执行 DNS 查询
	dnsStart := time.Now()
	dnsResult := utils.ResolveDNS(req.DomainName)
	dnsCost := time.Since(dnsStart).Milliseconds()

	if len(dnsResult.Results) == 0 {
		log.Warn("DNS 查询无结果",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.String("domain", req.DomainName),
		)
		response.ErrorWithTrace(c, 404, "没有找到相关 DNS 记录", traceID)
		return
	}

	// 构造响应
	resp := gin.H{
		"domain":     req.DomainName,
		"results":    dnsResult.Results,
		"count":      len(dnsResult.Results),
		"client_ip":  clientIP,
		"timestamp":  time.Now().UnixMilli(),
		"trace_id":   traceID,
		"query_cost": dnsCost,
	}

	log.Info("DNS 查询成功",
		zap.String("trace_id", traceID),
		zap.String("client_ip", clientIP),
		zap.String("domain", req.DomainName),
		zap.Int("result_count", len(dnsResult.Results)),
		zap.Int64("query_cost_ms", dnsCost),
	)

	response.Success(c, resp)
}

// RandomDomain 随机生成域名
func RandomDomain(c *gin.Context) {
	if IsPaused() {
		response.Error(c, 503, "系统已暂停，无法生成域名")
		return
	}

	domainName := GeneratingDomain()

	response.Success(c, gin.H{
		"domain": domainName,
	})
}

// ChangeServer 修改DNS服务器
func ChangeServer(c *gin.Context) {
	var dnsRequest ChangeDNSRequest
	if err := c.ShouldBindJSON(&dnsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Error("Failed to bind JSON", zap.Error(err))
		return
	}

	switch dnsRequest.Num {
	case 0, 1, 2:
		dns_server.ChangeServer(byte(dnsRequest.Num))
		server := dns_server.GetDNSServer(dnsRequest.Num)
		c.JSON(http.StatusOK, gin.H{"message": "DNS 服务器已更改为 " + server})
		log.Info("DNS 服务器已更改为", zap.String("server", server))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的选择"})
		log.Error("无效的选择", zap.Int("num", dnsRequest.Num))
	}
}

// ChangePact 修改协议
func ChangePact(c *gin.Context) {
	var pactRequest ChangePactRequest
	if err := c.ShouldBindJSON(&pactRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Error("Failed to bind JSON", zap.Error(err))
		return
	}

	switch pactRequest.Pact {
	case "udp":
		config.GlobalPact = "udp"
		c.JSON(http.StatusOK, gin.H{"message": "协议已更改为 UDP"})
		log.Info("协议已更改为 UDP")
	case "tcp":
		config.GlobalPact = "tcp"
		c.JSON(http.StatusOK, gin.H{"message": "协议已更改为 TCP"})
		log.Info("协议已更改为 TCP")
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的选择"})
	}
}
