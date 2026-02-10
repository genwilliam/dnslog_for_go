package dnslog

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type APIKey struct {
	ID         int64
	Name       string
	APIKey     string
	Enabled    bool
	CreatedAt  int64
	LastUsedAt int64
	Comment    string
}

type AuditLog struct {
	TraceID    string
	APIKeyID   sql.NullInt64
	Path       string
	Method     string
	ClientIP   string
	StatusCode int
	LatencyMs  int64
	Token      string
	CreatedAt  int64
}

var ErrAPIKeyNotFound = errors.New("api_key_not_found")
var ErrBootstrapConflict = errors.New("bootstrap_key_conflict")

func GetAPIKeyByHashWithContext(ctx context.Context, hash string) (APIKey, error) {
	if db == nil {
		return APIKey{}, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var k APIKey
	err := db.QueryRowContext(ctx, `
SELECT id, name, api_key, enabled, created_at, last_used_at, comment
FROM api_keys
WHERE api_key = ?
`, hash).Scan(&k.ID, &k.Name, &k.APIKey, &k.Enabled, &k.CreatedAt, &k.LastUsedAt, &k.Comment)
	if errors.Is(err, sql.ErrNoRows) {
		return APIKey{}, ErrAPIKeyNotFound
	}
	return k, err
}

func GetAPIKeyByHash(hash string) (APIKey, error) {
	return GetAPIKeyByHashWithContext(context.Background(), hash)
}

func TouchAPIKeyLastUsedWithContext(ctx context.Context, id int64, nowMs int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `
UPDATE api_keys SET last_used_at = ? WHERE id = ?
`, nowMs, id)
	return err
}

func TouchAPIKeyLastUsed(id int64, nowMs int64) error {
	return TouchAPIKeyLastUsedWithContext(context.Background(), id, nowMs)
}

func AddAuditLogWithContext(ctx context.Context, log AuditLog) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `
INSERT INTO audit_logs (trace_id, api_key_id, path, method, client_ip, status_code, latency_ms, token, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`, log.TraceID, log.APIKeyID, log.Path, log.Method, log.ClientIP, log.StatusCode, log.LatencyMs, log.Token, log.CreatedAt)
	return err
}

func AddAuditLog(log AuditLog) error {
	return AddAuditLogWithContext(context.Background(), log)
}

func CreateAPIKeyWithContext(ctx context.Context, name, hash, comment string, nowMs int64) (int64, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, `
INSERT INTO api_keys (name, api_key, enabled, created_at, last_used_at, comment)
VALUES (?, ?, 1, ?, 0, ?)
`, name, hash, nowMs, comment)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func CreateAPIKey(name, hash, comment string, nowMs int64) (int64, error) {
	return CreateAPIKeyWithContext(context.Background(), name, hash, comment, nowMs)
}

func SetAPIKeyEnabledWithContext(ctx context.Context, id int64, enabled bool) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	val := 0
	if enabled {
		val = 1
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, `
UPDATE api_keys SET enabled = ? WHERE id = ?
`, val, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}

func SetAPIKeyEnabled(id int64, enabled bool) error {
	return SetAPIKeyEnabledWithContext(context.Background(), id, enabled)
}

func ListAPIKeysWithContext(ctx context.Context, page, pageSize int) ([]APIKey, int, error) {
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
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM api_keys`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := db.QueryContext(ctx, `
SELECT id, name, api_key, enabled, created_at, last_used_at, comment
FROM api_keys
ORDER BY id DESC
LIMIT ? OFFSET ?
`, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []APIKey
	for rows.Next() {
		var k APIKey
		var enabled int
		if err := rows.Scan(&k.ID, &k.Name, &k.APIKey, &enabled, &k.CreatedAt, &k.LastUsedAt, &k.Comment); err != nil {
			return nil, 0, err
		}
		k.Enabled = enabled == 1
		items = append(items, k)
	}
	return items, total, nil
}

func ListAPIKeys(page, pageSize int) ([]APIKey, int, error) {
	return ListAPIKeysWithContext(context.Background(), page, pageSize)
}

func HasAPIKeysWithContext(ctx context.Context) (bool, error) {
	if db == nil {
		return false, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var exists int
	err := db.QueryRowContext(ctx, `SELECT 1 FROM api_keys WHERE enabled = 1 LIMIT 1`).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func HasAPIKeys() (bool, error) {
	return HasAPIKeysWithContext(context.Background())
}

func CountEnabledAPIKeysWithContext(ctx context.Context) (int, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM api_keys WHERE enabled = 1`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func CountEnabledAPIKeys() (int, error) {
	return CountEnabledAPIKeysWithContext(context.Background())
}

func countEnabledAPIKeysForUpdate(ctx context.Context, tx *sql.Tx) (int, error) {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM api_keys WHERE enabled = 1 FOR UPDATE`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func CreateBootstrapAPIKeyWithContext(ctx context.Context, name, hash, comment string, nowMs int64) (int64, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	count, err := countEnabledAPIKeysForUpdate(ctx, tx)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, ErrBootstrapConflict
	}

	res, err := tx.ExecContext(ctx, `
INSERT INTO api_keys (name, api_key, enabled, created_at, last_used_at, comment)
VALUES (?, ?, 1, ?, 0, ?)
`, name, hash, nowMs, comment)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func CreateBootstrapAPIKey(name, hash, comment string, nowMs int64) (int64, error) {
	return CreateBootstrapAPIKeyWithContext(context.Background(), name, hash, comment, nowMs)
}
