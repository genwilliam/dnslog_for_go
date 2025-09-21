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
	var domain struct {
		DomainName string `json:"domain_name"` // 接收域名
	}

	// 解析JSON数据
	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查域名是否合法
	if !utils.StandardizeDomain(domain.DomainName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名不合法，请重新输入"})
		return
	}

	// DNS 查询
	dnsResult := utils.ResolveDNS(domain.DomainName)

	// 返回查询结果
	if len(dnsResult.Results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到相关的 DNS 记录"})
		log.Error("没有找到相关的 DNS 记录", zap.String("domain", domain.DomainName))
	} else {
		c.JSON(http.StatusOK, dnsResult)
	}
}

// RandomDomain 随机生成域名
func RandomDomain(c *gin.Context) {
	domainName := GeneratingDomain()
	c.JSON(http.StatusOK, gin.H{"domain": domainName})
}

// ChangeServer 修改DNS服务器
func ChangeServer(c *gin.Context) {
	var dnsRequest ChangeDNSRequest

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&dnsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Error("Failed to bind JSON", zap.Error(err))
		return
	}

	switch dnsRequest.Num {
	case 0:
		c.JSON(http.StatusOK, gin.H{"message": "DNS 服务器已更改为 8.8.8.8 (Google)"})
		log.Info("DNS 服务器已更改为 8.8.8.8 (Google)")
		dns_server.ChangeServer(0)
	case 1:
		c.JSON(http.StatusOK, gin.H{"message": "DNS 服务器已更改为 223.5.5.5 (阿里)"})
		log.Info("DNS 服务器已更改为 223.5.5.5 (阿里)")
		dns_server.ChangeServer(1)
	case 2:
		c.JSON(http.StatusOK, gin.H{"message": "DNS 服务器已更改为 127.0.0.1(本地)"})
		log.Info("DNS 服务器已更改为 127.0.0.1 (本地)")
		dns_server.ChangeServer(2)
	default:
		log.Error("无效的选择", zap.Int("num", dnsRequest.Num))
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的选择"})
	}
}

// ChangePact 修改协议
func ChangePact(c *gin.Context) {
	var pactRequest ChangePactRequest

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&pactRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		log.Error("Failed to bind JSON", zap.Error(err))
		return
	}

	// 根据 pactRequest.Pact 的值更新全局协议
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
