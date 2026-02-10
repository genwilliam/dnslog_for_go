package dnslog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type TokenStatus struct {
	Token     string `json:"token"`
	Domain    string `json:"domain"`
	Status    string `json:"status"`
	FirstSeen int64  `json:"first_seen"`
	LastSeen  int64  `json:"last_seen"`
	HitCount  int64  `json:"hit_count"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	ExpiresAt int64  `json:"expires_at"`
}

type TokenListFilter struct {
	Page        int
	PageSize    int
	Status      string
	CreatedStart int64
	CreatedEnd   int64
	LastStart    int64
	LastEnd      int64
	Keyword     string
	OrderBy     string
	Order       string
}

var ErrTokenNotFound = errors.New("token_not_found")

func CreateTokenInitWithContext(ctx context.Context, token, domain string, nowMs, expiresAtMs int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
INSERT INTO dns_tokens (token, domain, status, hit_count, first_seen, last_seen, created_at, updated_at, expires_at)
VALUES (?, ?, 'INIT', 0, 0, 0, ?, ?, ?)
`, token, domain, nowMs, nowMs, expiresAtMs)
	if isDuplicateKey(err) {
		return fmt.Errorf("token exists: %w", err)
	}
	return err
}

func CreateTokenInit(token, domain string, nowMs, expiresAtMs int64) error {
	return CreateTokenInitWithContext(context.Background(), token, domain, nowMs, expiresAtMs)
}

func UpsertTokenHitWithContext(ctx context.Context, token, domain string, nowMs, ttlMs int64) (bool, error) {
	if db == nil {
		return false, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	expiresAt := nowMs + ttlMs
	res, err := db.ExecContext(ctx, `
INSERT INTO dns_tokens (token, domain, status, hit_count, first_seen, last_seen, created_at, updated_at, expires_at)
VALUES (?, ?, 'HIT', 1, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  hit_count = IF(status = 'EXPIRED', hit_count + LAST_INSERT_ID(0), LAST_INSERT_ID(hit_count + 1)),
  last_seen = IF(status = 'EXPIRED', last_seen, VALUES(last_seen)),
  updated_at = VALUES(updated_at),
  status = IF(status = 'EXPIRED', 'EXPIRED', 'HIT'),
  first_seen = IF(first_seen = 0, VALUES(first_seen), first_seen),
  expires_at = IF(status = 'EXPIRED', expires_at, VALUES(expires_at))
`, token, domain, nowMs, nowMs, nowMs, nowMs, expiresAt)
	if err != nil {
		return false, err
	}

	affected, _ := res.RowsAffected()
	if affected == 1 {
		return true, nil
	}
	if affected == 0 {
		return false, nil
	}
	hitCount, err := res.LastInsertId()
	if err != nil {
		return false, nil
	}
	return hitCount == 1, nil
}

func UpsertTokenHit(token, domain string, nowMs, ttlMs int64) (bool, error) {
	return UpsertTokenHitWithContext(context.Background(), token, domain, nowMs, ttlMs)
}

func GetTokenStatusWithContext(ctx context.Context, token string) (TokenStatus, error) {
	if db == nil {
		return TokenStatus{}, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var ts TokenStatus
	err := db.QueryRowContext(ctx, `
SELECT token, domain, status, first_seen, last_seen, hit_count, created_at, updated_at, expires_at
FROM dns_tokens
WHERE token = ?
`, token).Scan(&ts.Token, &ts.Domain, &ts.Status, &ts.FirstSeen, &ts.LastSeen, &ts.HitCount, &ts.CreatedAt, &ts.UpdatedAt, &ts.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return TokenStatus{}, ErrTokenNotFound
	}
	return ts, err
}

func GetTokenStatus(token string) (TokenStatus, error) {
	return GetTokenStatusWithContext(context.Background(), token)
}

func MaybeExpireTokenWithContext(ctx context.Context, token string, nowMs int64) (bool, error) {
	if db == nil {
		return false, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, `
UPDATE dns_tokens
SET status = 'EXPIRED', updated_at = ?
WHERE token = ? AND status != 'EXPIRED' AND expires_at <= ?
`, nowMs, token, nowMs)
	if err != nil {
		return false, err
	}
	affected, _ := res.RowsAffected()
	return affected > 0, nil
}

func MaybeExpireToken(token string, nowMs int64) (bool, error) {
	return MaybeExpireTokenWithContext(context.Background(), token, nowMs)
}

func MarkExpiredBatchWithContext(ctx context.Context, nowMs int64, limit int) (int64, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	if limit <= 0 {
		limit = 200
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := db.ExecContext(ctx, `
UPDATE dns_tokens
SET status = 'EXPIRED', updated_at = ?
WHERE status != 'EXPIRED' AND expires_at <= ?
ORDER BY expires_at ASC
LIMIT ?
`, nowMs, nowMs, limit)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, nil
}

func MarkExpiredBatch(nowMs int64, limit int) (int64, error) {
	return MarkExpiredBatchWithContext(context.Background(), nowMs, limit)
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

func IsDuplicateKeyError(err error) bool {
	return isDuplicateKey(err)
}

func ListTokensWithContext(ctx context.Context, filter TokenListFilter) ([]TokenStatus, int, error) {
	if db == nil {
		return nil, 0, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	order := "DESC"
	if strings.ToLower(filter.Order) == "asc" {
		order = "ASC"
	}
	orderBy := "created_at"
	if strings.ToLower(filter.OrderBy) == "last_seen" {
		orderBy = "last_seen"
	}

	where := make([]string, 0)
	args := make([]interface{}, 0)
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Keyword != "" {
		where = append(where, "(token LIKE ? OR domain LIKE ?)")
		kw := "%" + filter.Keyword + "%"
		args = append(args, kw, kw)
	}
	if filter.CreatedStart > 0 {
		where = append(where, "created_at >= ?")
		args = append(args, filter.CreatedStart)
	}
	if filter.CreatedEnd > 0 {
		where = append(where, "created_at <= ?")
		args = append(args, filter.CreatedEnd)
	}
	if filter.LastStart > 0 {
		where = append(where, "last_seen >= ?")
		args = append(args, filter.LastStart)
	}
	if filter.LastEnd > 0 {
		where = append(where, "last_seen <= ?")
		args = append(args, filter.LastEnd)
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	countSQL := "SELECT COUNT(1) FROM dns_tokens " + whereSQL
	var total int
	if err := db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tokens: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	querySQL := `
SELECT token, domain, status, first_seen, last_seen, hit_count, created_at, updated_at, expires_at
FROM dns_tokens
` + whereSQL + `
ORDER BY ` + orderBy + ` ` + order + `
LIMIT ? OFFSET ?`

	argsWithPage := append(args, filter.PageSize, offset)
	rows, err := db.QueryContext(ctx, querySQL, argsWithPage...)
	if err != nil {
		return nil, 0, fmt.Errorf("query tokens: %w", err)
	}
	defer rows.Close()

	var items []TokenStatus
	for rows.Next() {
		var ts TokenStatus
		if err := rows.Scan(&ts.Token, &ts.Domain, &ts.Status, &ts.FirstSeen, &ts.LastSeen, &ts.HitCount, &ts.CreatedAt, &ts.UpdatedAt, &ts.ExpiresAt); err != nil {
			return nil, 0, fmt.Errorf("scan token: %w", err)
		}
		items = append(items, ts)
	}

	return items, total, nil
}

func ListTokens(filter TokenListFilter) ([]TokenStatus, int, error) {
	return ListTokensWithContext(context.Background(), filter)
}
