package domain

import (
	"net/http"
	"strings"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"github.com/genwilliam/dnslog_for_go/pkg/response"
	"github.com/genwilliam/dnslog_for_go/pkg/utils"

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

	traceID := response.GetTraceID(c)
	if traceID == "" {
		traceID = utils.GenerateTraceID()
	}
	clientIP := c.ClientIP()

	if IsPaused() {
		log.Warn("系统暂停，拒绝查询",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
		)
		response.ErrorWithTrace(c, 503, response.CodeSystemPaused, traceID)
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
		response.ErrorWithTrace(c, 400, response.CodeBadRequest, traceID)
		return
	}

	if req.DomainName == "" {
		response.ErrorWithTrace(c, 400, response.CodeBadRequest, traceID)
		return
	}

	normalizedDomain := normalizeDomain(req.DomainName)
	if !utils.StandardizeDomain(normalizedDomain) {
		log.Info("域名格式不合法",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.String("domain", normalizedDomain),
		)
		response.ErrorWithTrace(c, 400, response.CodeBadRequest, traceID)
		return
	}

	cfg := config.Get()
	if !isDomainAllowed(normalizedDomain, cfg) {
		log.Warn("域名不在允许的根域范围内",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.String("domain", normalizedDomain),
		)
		response.ErrorWithTrace(c, 403, response.CodeForbidden, traceID)
		return
	}

	filter := dnslog.ListFilter{
		Page:     1,
		PageSize: cfg.MaxPageSize,
		Domain:   normalizedDomain,
	}
	items, total, err := dnslog.ListRecordsWithContext(c.Request.Context(), filter)
	if err != nil {
		log.Error("查询 DNS 记录失败",
			zap.String("trace_id", traceID),
			zap.String("client_ip", clientIP),
			zap.String("domain", normalizedDomain),
			zap.Error(err),
		)
		response.ErrorWithTrace(c, 500, response.CodeInternalError, traceID)
		return
	}

	resp := gin.H{
		"domain":    normalizedDomain,
		"items":     items,
		"total":     total,
		"pending":   total == 0,
		"timestamp": time.Now().UnixMilli(),
		"trace_id":  traceID,
	}

	log.Info("DNS 记录查询完成",
		zap.String("trace_id", traceID),
		zap.String("client_ip", clientIP),
		zap.String("domain", normalizedDomain),
		zap.Int("record_total", total),
	)

	response.Success(c, resp)
}

// RandomDomain 随机生成域名
func RandomDomain(c *gin.Context) {
	if IsPaused() {
		response.Error(c, 503, response.CodeSystemPaused)
		return
	}

	domainName, token, err := GenerateAndInitDomainWithContext(c.Request.Context())
	if err != nil {
		log.Error("生成域名失败", zap.Error(err))
		response.Error(c, 500, response.CodeInternalError)
		return
	}

	response.Success(c, gin.H{
		"domain": domainName,
		"token":  token,
	})
}

// ChangeServer 修改DNS服务器
func ChangeServer(c *gin.Context) {
	var dnsRequest ChangeDNSRequest
	if err := c.ShouldBindJSON(&dnsRequest); err != nil {
		log.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	cfg := config.Get()
	if dnsRequest.Num < 0 || dnsRequest.Num >= len(cfg.UpstreamDNS) {
		log.Error("无效的选择", zap.Int("num", dnsRequest.Num))
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	cfg.SetUpstreamIndex(dnsRequest.Num)
	server := cfg.CurrentUpstream()
	log.Info("DNS 服务器已更改为", zap.String("server", server))
	response.Success(c, gin.H{"message": "DNS 服务器已更改为 " + server})
}

// ChangePact 修改协议
func ChangePact(c *gin.Context) {
	var pactRequest ChangePactRequest
	if err := c.ShouldBindJSON(&pactRequest); err != nil {
		log.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
		return
	}

	cfg := config.Get()

	switch pactRequest.Pact {
	case "udp":
		cfg.SetProtocol("udp")
		log.Info("协议已更改为 UDP")
		response.Success(c, gin.H{"message": "协议已更改为 UDP"})
	case "tcp":
		cfg.SetProtocol("tcp")
		log.Info("协议已更改为 TCP")
		response.Success(c, gin.H{"message": "协议已更改为 TCP"})
	default:
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest)
	}
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))
	return strings.TrimSuffix(domain, ".")
}

func isDomainAllowed(domain string, cfg *config.Config) bool {
	if cfg.CaptureAll {
		return true
	}
	normalized := normalizeDomain(domain)
	if normalized == "" {
		return false
	}
	roots := make([]string, 0, 1+len(cfg.RootDomains))
	if cfg.RootDomain != "" {
		roots = append(roots, cfg.RootDomain)
	}
	roots = append(roots, cfg.RootDomains...)

	for _, root := range roots {
		r := normalizeDomain(root)
		if r == "" {
			continue
		}
		if normalized == r || strings.HasSuffix(normalized, "."+r) {
			return true
		}
	}
	return false
}
