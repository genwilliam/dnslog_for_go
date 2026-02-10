package dnslog

import (
	"context"
	"fmt"
	"time"

	"github.com/genwilliam/dnslog_for_go/internal/infra"
)

func AllowDNSQuery(clientIP string) bool {
	if activeConfig == nil || !activeConfig.DNSRateLimitEnabled {
		return true
	}
	client := infra.GetRedis()
	if client == nil {
		return true
	}
	key := fmt.Sprintf("dns_rl:%s", clientIP)
	limit := activeConfig.DNSRateLimitMaxRequests
	if limit <= 0 {
		limit = 1000
	}
	windowSeconds := activeConfig.DNSRateLimitWindowSeconds
	if windowSeconds <= 0 {
		windowSeconds = 60
	}
	window := time.Duration(windowSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return true
	}
	if count == 1 {
		_ = client.Expire(ctx, key, window).Err()
	}
	return int(count) <= limit
}
