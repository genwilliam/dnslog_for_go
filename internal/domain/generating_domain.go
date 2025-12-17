package domain

import (
	"fmt"
	"strings"

	"github.com/genwilliam/dnslog_for_go/config"
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

	root := strings.TrimSuffix(cfg.RootDomain, ".")
	domain := fmt.Sprintf("%s.%s", token, root)

	log.Info("生成域名", zap.String("token", token), zap.String("domain", domain))
	return domain
}
