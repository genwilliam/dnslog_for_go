package router

import (
	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
	"strings"
)

// ConfigHandler 返回当前运行时配置（非敏感字段）
func ConfigHandler(c *gin.Context) {
	cfg := config.Get()

	dnsPort := "53"
	if idx := strings.LastIndex(cfg.DNSListenAddr, ":"); idx >= 0 && idx < len(cfg.DNSListenAddr)-1 {
		dnsPort = cfg.DNSListenAddr[idx+1:]
	}

	response.Success(c, gin.H{
		"root_domain":       cfg.RootDomain,
		"dns_listen_addr":   cfg.DNSListenAddr,
		"http_listen":       cfg.HTTPListenAddr,
		"upstream_dns":      cfg.UpstreamDNS,
		"protocol":          cfg.GetProtocol(),
		"page_size":         cfg.DefaultPageSize,
		"max_page_size":     cfg.MaxPageSize,
		"token_ttl":         cfg.TokenTTLSeconds,
		"apiKeyRequired":    cfg.APIKeyRequired,   // camelCase 兼容
		"api_key_required":  cfg.APIKeyRequired,   // snake_case 兼容
		"dns_port":          dnsPort,
	})
}
