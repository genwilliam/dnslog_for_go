package dnslog

import (
	"context"
	"errors"
	"time"
)

func DeleteOldRecords(cutoffMs int64, limit int) (int64, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	if limit <= 0 {
		limit = 1000
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, `
DELETE FROM dns_records
WHERE timestamp < ?
LIMIT ?
`, cutoffMs, limit)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, nil
}
