package dnslog

import (
	"time"

	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"go.uber.org/zap"
)

func StartExpireWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			nowMs := time.Now().UnixMilli()
			affected, err := MarkExpiredBatch(nowMs, 500)
			if err != nil {
				log.Error("expire worker failed", zap.Error(err))
				continue
			}
			if affected > 0 {
				log.Info("expired tokens updated", zap.Int64("count", affected))
			}
		}
	}()
}
