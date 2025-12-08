package utils

import (
	"regexp"
	"strings"

	"github.com/genwilliam/dnslog_for_go/pkg/log"
)

// 允许的常见 TLD 后缀（可扩展）
var commonTLDs = []string{
	".com", ".net", ".org", ".cn", ".io", ".edu", ".gov", ".co", ".xyz",
}

// StandardizeDomain 判断域名是否合法
func StandardizeDomain(domain string) bool {
	// 转成小写
	domain = strings.ToLower(domain)

	// 长度限制
	if len(domain) < 3 || len(domain) > 253 {
		log.Warn("域名太短或太长，长度应在 3 到 253 个字符之间")
		return false
	}

	// 检查是否有合法的 TLD
	validTLD := false
	for _, tld := range commonTLDs {
		if strings.HasSuffix(domain, tld) {
			validTLD = true
			break
		}
	}
	if !validTLD {
		log.Warn("域名必须以常见的顶级域名结尾，例如 .com、.net")
		return false
	}

	// 正则校验整个域名结构是否合法
	match, _ := regexp.MatchString(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)(?:\.(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?))*\.[a-z]{2,}$`, domain)
	if !match {
		log.Warn("域名格式不正确")
		return false
	}

	log.Info("域名格式正确")
	return true
}
