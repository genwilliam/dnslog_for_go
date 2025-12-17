package domain_test

import (
	"strings"
	"testing"

	"github.com/genwilliam/dnslog_for_go/config"
)

// 基于 uuid 生成的域名应落在配置的根域名下。
func TestGeneratingDomain(t *testing.T) {
	config.Load()

	domain := GeneratingDomain()
	if !strings.Contains(domain, ".") {
		t.Fatalf("生成域名不合法: %s", domain)
	}

	root := config.Get().RootDomain
	if !strings.HasSuffix(domain, root) {
		t.Fatalf("生成域名未使用根域名, got=%s root=%s", domain, root)
	}
}
