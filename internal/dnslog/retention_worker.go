package dnslog

import (
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"go.uber.org/zap"
)

func StartRetentionWorker(cfg *config.Config) {
	if cfg == nil || !cfg.RetentionEnabled {
		return
	}
	if cfg.RecordRetentionDays <= 0 {
		return
	}
	interval := time.Duration(cfg.RetentionIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = time.Hour
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-time.Duration(cfg.RecordRetentionDays) * 24 * time.Hour).UnixMilli()
			affected, err := DeleteOldRecords(cutoff, cfg.RetentionBatchSize)
			if err != nil {
				log.Error("retention cleanup failed", zap.Error(err))
				continue
			}
			if affected > 0 {
				log.Info("retention cleanup", zap.Int64("deleted", affected))
			}
		}
	}()
}
