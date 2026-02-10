package dnslog

import (
	"context"
	"errors"
	"time"

	"github.com/genwilliam/dnslog_for_go/internal/infra"
)

type BlacklistEntry struct {
	ID        int64  `json:"id"`
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at"`
}

var ErrBlacklistNotFound = errors.New("blacklist_not_found")

const blacklistRedisKey = "blacklist:ip"

func IsIPBlacklistedWithContext(ctx context.Context, ip string) (bool, error) {
	if db == nil {
		return false, errors.New("store not initialized")
	}
	if ip == "" {
		return false, nil
	}

	if client := infra.GetRedis(); client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		ok, err := client.SIsMember(ctx, blacklistRedisKey, ip).Result()
		if err == nil {
			return ok, nil
		}
	}

	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var enabled int
	err := db.QueryRowContext(ctx, `
SELECT enabled FROM ip_blacklist WHERE ip = ?
`, ip).Scan(&enabled)
	if err != nil {
		return false, nil
	}
	if enabled == 1 {
		_ = addIPToRedis(ip)
		return true, nil
	}
	return false, nil
}

func IsIPBlacklisted(ip string) (bool, error) {
	return IsIPBlacklistedWithContext(context.Background(), ip)
}

func AddBlacklistIPWithContext(ctx context.Context, ip, reason string, nowMs int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `
INSERT INTO ip_blacklist (ip, reason, enabled, created_at)
VALUES (?, ?, 1, ?)
ON DUPLICATE KEY UPDATE enabled = 1, reason = VALUES(reason)
`, ip, reason, nowMs)
	if err == nil {
		_ = addIPToRedis(ip)
	}
	return err
}

func AddBlacklistIP(ip, reason string, nowMs int64) error {
	return AddBlacklistIPWithContext(context.Background(), ip, reason, nowMs)
}

func DisableBlacklistIPWithContext(ctx context.Context, id int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var ip string
	if err := db.QueryRowContext(ctx, `SELECT ip FROM ip_blacklist WHERE id = ?`, id).Scan(&ip); err != nil {
		return ErrBlacklistNotFound
	}
	if _, err := db.ExecContext(ctx, `
UPDATE ip_blacklist SET enabled = 0 WHERE id = ?
`, id); err != nil {
		return err
	}
	_ = removeIPFromRedis(ip)
	return nil
}

func DisableBlacklistIP(id int64) error {
	return DisableBlacklistIPWithContext(context.Background(), id)
}

func ListBlacklistWithContext(ctx context.Context, page, pageSize int) ([]BlacklistEntry, int, error) {
	if db == nil {
		return nil, 0, errors.New("store not initialized")
	}
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var total int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM ip_blacklist`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := db.QueryContext(ctx, `
SELECT id, ip, reason, enabled, created_at
FROM ip_blacklist
ORDER BY id DESC
LIMIT ? OFFSET ?
`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []BlacklistEntry
	for rows.Next() {
		var e BlacklistEntry
		var enabled int
		if err := rows.Scan(&e.ID, &e.IP, &e.Reason, &enabled, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		e.Enabled = enabled == 1
		items = append(items, e)
	}
	return items, total, nil
}

func ListBlacklist(page, pageSize int) ([]BlacklistEntry, int, error) {
	return ListBlacklistWithContext(context.Background(), page, pageSize)
}

func addIPToRedis(ip string) error {
	client := infra.GetRedis()
	if client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	return client.SAdd(ctx, blacklistRedisKey, ip).Err()
}

func removeIPFromRedis(ip string) error {
	client := infra.GetRedis()
	if client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	return client.SRem(ctx, blacklistRedisKey, ip).Err()
}
