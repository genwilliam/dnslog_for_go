package utils

import (
	"regexp"
	"strings"

	"github.com/genwilliam/dnslog_for_go/pkg/log"
)

// StandardizeDomain 判断域名是否合法
func StandardizeDomain(domain string) bool {
	// 转成小写
	domain = strings.ToLower(domain)

	// 长度限制
	if len(domain) < 3 || len(domain) > 253 {
		log.Warn("域名太短或太长，长度应在 3 到 253 个字符之间")
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
