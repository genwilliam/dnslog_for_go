package router

import (
	_ "net/http"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
)

// ConfigHandler 返回当前运行时配置（非敏感字段）
func ConfigHandler(c *gin.Context) {
	cfg := config.Get()
	response.Success(c, gin.H{
		"root_domain":     cfg.RootDomain,
		"dns_listen_addr": cfg.DNSListenAddr,
		"http_listen":     cfg.HTTPListenAddr,
		"upstream_dns":    cfg.UpstreamDNS,
		"protocol":        cfg.Protocol,
		"page_size":       cfg.DefaultPageSize,
		"max_page_size":   cfg.MaxPageSize,
	})
}
