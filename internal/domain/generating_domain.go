package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GeneratingDomain 基于 uuid 生成子域名，并拼接根域。
func GeneratingDomain() string {
	cfg := config.Get()

	id := uuid.New().String()
	cleaned := strings.ReplaceAll(id, "-", "")

	if len(cleaned) < 10 {
		return fmt.Sprintf("uuid too short: %s", cleaned)
	}
	token := cleaned[:10]

	root := selectRootDomain(cfg)
	domain := fmt.Sprintf("%s.%s", token, root)

	log.Info("生成域名", zap.String("token", token), zap.String("domain", domain))
	return domain
}

// GenerateAndInitDomain 生成域名并写入 token 状态表（INIT）
func GenerateAndInitDomainWithContext(ctx context.Context) (string, string, error) {
	cfg := config.Get()
	nowMs := time.Now().UnixMilli()
	ttlMs := int64(cfg.TokenTTLSeconds) * 1000
	if ttlMs <= 0 {
		ttlMs = int64(3600 * 1000)
	}

	root := selectRootDomain(cfg)
	if root == "" {
		return "", "", fmt.Errorf("root domain is empty")
	}

	for i := 0; i < 5; i++ {
		id := uuid.New().String()
		cleaned := strings.ReplaceAll(id, "-", "")
		if len(cleaned) < 10 {
			continue
		}
		token := cleaned[:10]
		domain := fmt.Sprintf("%s.%s", token, root)
		expiresAt := nowMs + ttlMs

		if err := dnslog.CreateTokenInitWithContext(ctx, token, domain, nowMs, expiresAt); err != nil {
			if dnslog.IsDuplicateKeyError(err) {
				log.Warn("token 冲突，准备重试", zap.Error(err))
				continue
			}
			return "", "", err
		}
		log.Info("生成并初始化 token", zap.String("token", token), zap.String("domain", domain))
		return domain, token, nil
	}
	return "", "", fmt.Errorf("failed to allocate token")
}

// GenerateAndInitDomain 生成域名并写入 token 状态表（INIT）
func GenerateAndInitDomain() (string, string, error) {
	return GenerateAndInitDomainWithContext(context.Background())
}

func selectRootDomain(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	if cfg.RootDomain != "" {
		return strings.TrimSuffix(cfg.RootDomain, ".")
	}
	if len(cfg.RootDomains) > 0 {
		return strings.TrimSuffix(cfg.RootDomains[0], ".")
	}
	return ""
}
