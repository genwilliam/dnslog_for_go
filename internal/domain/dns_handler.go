package domain

import (
	"dnslog_for_go/internal/config"
	"dnslog_for_go/internal/domain/dns_server"
	"dnslog_for_go/pkg/log"
	"dnslog_for_go/pkg/utils"
	"net/http"

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
	if IsPaused() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "系统已暂停，无法查询域名"})
		log.Warn("系统已暂停，无法查询域名")
		return
	}

	var domain struct {
		DomainName string `json:"domain_name"`
	}

	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.StandardizeDomain(domain.DomainName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名不合法，请重新输入"})
		return
	}

	dnsResult := utils.ResolveDNS(domain.DomainName)
	if len(dnsResult.Results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到相关的 DNS 记录"})
		log.Error("没有找到相关的 DNS 记录", zap.String("domain", domain.DomainName))
	} else {
		c.JSON(http.StatusOK, dnsResult)
	}
}

// RandomDomain 随机生成域名
func RandomDomain(c *gin.Context) {
	if IsPaused() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "系统已暂停，无法生成域名"})
		log.Warn("系统已暂停，无法生成域名")
		return
	}

	domainName := GeneratingDomain()
	c.JSON(http.StatusOK, gin.H{"domain": domainName})
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
