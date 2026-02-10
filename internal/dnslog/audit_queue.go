package dnslog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/genwilliam/dnslog_for_go/internal/infra"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"go.uber.org/zap"
)

const auditQueueKey = "audit:queue"

func EnqueueAuditLog(logEntry AuditLog) error {
	client := infra.GetRedis()
	if client == nil {
		return AddAuditLog(logEntry)
	}
	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return client.LPush(ctx, auditQueueKey, data).Err()
}

func StartAuditWorker() {
	client := infra.GetRedis()
	if client == nil {
		return
	}
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := client.BRPop(ctx, 3*time.Second, auditQueueKey).Result()
			cancel()
			if err != nil || len(res) < 2 {
				continue
			}
			var entry AuditLog
			if err := json.Unmarshal([]byte(res[1]), &entry); err != nil {
				log.Error("decode audit log failed", zap.Error(err))
				continue
			}
			if err := AddAuditLog(entry); err != nil {
				log.Error("write audit log failed", zap.Error(err))
			}
		}
	}()
}
