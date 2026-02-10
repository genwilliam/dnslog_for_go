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

// Record 表示一条 DNS 请求日志
type Record struct {
	ID        int64  `json:"id"`
	Domain    string `json:"domain"`
	ClientIP  string `json:"client_ip"`
	Protocol  string `json:"protocol"`
	QType     string `json:"qtype"`
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Token     string `json:"token"`
}

// ListFilter 查询过滤条件
type ListFilter struct {
	Page     int
	PageSize int

	Domain   string
	ClientIP string
	Protocol string
	QType    string
	Token    string
	Order    string
	Cursor   int64

	Start int64 // 起始时间戳（毫秒）
	End   int64 // 结束时间戳（毫秒）
}

var db *sql.DB

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// InitStore 初始化 MySQL 连接与表结构。
func InitStore(dsn string) error {
	if db != nil {
		return nil
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("invalid mysql dsn: %w", err)
	}

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("open mysql: %w", err)
	}

	conn.SetConnMaxLifetime(10 * time.Minute)
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(20)

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping mysql: %w", err)
	}

	if err := createTable(conn); err != nil {
		return err
	}
	if err := createTokensTable(conn); err != nil {
		return err
	}
	if err := createAPIKeysTable(conn); err != nil {
		return err
	}
	if err := createAuditLogsTable(conn); err != nil {
		return err
	}
	if err := createIPBlacklistTable(conn); err != nil {
		return err
	}
	if err := createTokenWebhooksTable(conn); err != nil {
		return err
	}
	if err := createWebhookJobsTable(conn); err != nil {
		return err
	}

	db = conn
	return nil
}

func createTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS dns_records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    domain     VARCHAR(255) NOT NULL,
    client_ip  VARCHAR(64)  DEFAULT '',
    protocol   VARCHAR(16)  DEFAULT '',
    qtype      VARCHAR(32)  DEFAULT '',
    timestamp  BIGINT       NOT NULL,
    server     VARCHAR(64)  DEFAULT '',
    token      VARCHAR(128) DEFAULT '',
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_domain (domain),
    INDEX idx_token (token),
    INDEX idx_ts (timestamp),
    INDEX idx_token_ts (token, timestamp),
    INDEX idx_qtype_ts (qtype, timestamp),
    INDEX idx_client_ts (client_ip, timestamp),
    INDEX idx_protocol_ts (protocol, timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table dns_records: %w", err)
	}
	return nil
}

func createTokensTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS dns_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token       VARCHAR(128) NOT NULL UNIQUE,
    domain      VARCHAR(255) NOT NULL,
    status      ENUM('INIT','HIT','EXPIRED') NOT NULL DEFAULT 'INIT',
    hit_count   BIGINT NOT NULL DEFAULT 0,
    first_seen  BIGINT NOT NULL DEFAULT 0,
    last_seen   BIGINT NOT NULL DEFAULT 0,
    created_at  BIGINT NOT NULL,
    updated_at  BIGINT NOT NULL,
    expires_at  BIGINT NOT NULL,
    INDEX idx_status (status),
    INDEX idx_expires (expires_at),
    INDEX idx_status_created (status, created_at),
    INDEX idx_status_last (status, last_seen)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table dns_tokens: %w", err)
	}
	return nil
}

func createAPIKeysTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    api_key VARCHAR(128) NOT NULL UNIQUE,
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    last_used_at BIGINT NOT NULL DEFAULT 0,
    comment VARCHAR(255) DEFAULT '',
    INDEX idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table api_keys: %w", err)
	}
	return nil
}

func createAuditLogsTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL,
    api_key_id BIGINT DEFAULT NULL,
    path VARCHAR(128) NOT NULL,
    method VARCHAR(16) NOT NULL,
    client_ip VARCHAR(64) NOT NULL,
    status_code INT NOT NULL,
    latency_ms INT NOT NULL,
    token VARCHAR(128) DEFAULT '',
    created_at BIGINT NOT NULL,
    INDEX idx_created (created_at),
    INDEX idx_path (path),
    INDEX idx_ip (client_ip)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table audit_logs: %w", err)
	}
	return nil
}

func createIPBlacklistTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS ip_blacklist (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ip VARCHAR(64) NOT NULL UNIQUE,
    reason VARCHAR(255) DEFAULT '',
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    INDEX idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table ip_blacklist: %w", err)
	}
	return nil
}

func createTokenWebhooksTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS token_webhooks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(128) NOT NULL,
    webhook_url VARCHAR(512) NOT NULL,
    secret VARCHAR(128) DEFAULT '',
    mode ENUM('FIRST_HIT','EACH_HIT') NOT NULL DEFAULT 'FIRST_HIT',
    enabled TINYINT NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    UNIQUE KEY uk_token (token),
    INDEX idx_token (token),
    INDEX idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table token_webhooks: %w", err)
	}
	return nil
}

func createWebhookJobsTable(conn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS webhook_jobs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(128) NOT NULL,
    url VARCHAR(512) NOT NULL,
    payload TEXT NOT NULL,
    secret VARCHAR(128) DEFAULT '',
    status ENUM('PENDING','SUCCESS','FAILED') NOT NULL DEFAULT 'PENDING',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at BIGINT NOT NULL DEFAULT 0,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_status_retry (status, next_retry_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table webhook_jobs: %w", err)
	}
	return nil
}

// AddRecord 持久化一条记录。
func AddRecordWithContext(ctx context.Context, rec Record) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
INSERT INTO dns_records (domain, client_ip, protocol, qtype, timestamp, server, token)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.Domain, rec.ClientIP, rec.Protocol, rec.QType, rec.Timestamp, rec.Server, rec.Token)
	return err
}

// AddRecord 持久化一条记录。
func AddRecord(rec Record) error {
	return AddRecordWithContext(context.Background(), rec)
}

// ListRecords 根据过滤条件分页查询。
func ListRecordsWithContext(ctx context.Context, filter ListFilter) ([]Record, int, error) {
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

	where := make([]string, 0)
	args := make([]interface{}, 0)

	addLike := func(col, val string) {
		if val != "" {
			where = append(where, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+val+"%")
		}
	}

	addEq := func(col, val string) {
		if val != "" {
			where = append(where, fmt.Sprintf("%s = ?", col))
			args = append(args, val)
		}
	}

	addLike("domain", filter.Domain)
	addLike("client_ip", filter.ClientIP)
	addLike("token", filter.Token)
	addEq("protocol", filter.Protocol)
	addEq("qtype", filter.QType)

	if filter.Start > 0 {
		where = append(where, "timestamp >= ?")
		args = append(args, filter.Start)
	}
	if filter.End > 0 {
		where = append(where, "timestamp <= ?")
		args = append(args, filter.End)
	}
	if filter.Cursor > 0 {
		if order == "ASC" {
			where = append(where, "timestamp > ?")
		} else {
			where = append(where, "timestamp < ?")
		}
		args = append(args, filter.Cursor)
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	countSQL := "SELECT COUNT(1) FROM dns_records " + whereSQL
	var total int
	if err := db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	if filter.Cursor > 0 {
		offset = 0
	}
	querySQL := `
SELECT id, domain, client_ip, protocol, qtype, timestamp, server, token
FROM dns_records
` + whereSQL + `
ORDER BY timestamp ` + order + `
LIMIT ? OFFSET ?`

	argsWithPage := append(args, filter.PageSize, offset)

	rows, err := db.QueryContext(ctx, querySQL, argsWithPage...)
	if err != nil {
		return nil, 0, fmt.Errorf("query records: %w", err)
	}
	defer rows.Close()

	var items []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.Domain, &rec.ClientIP, &rec.Protocol, &rec.QType, &rec.Timestamp, &rec.Server, &rec.Token); err != nil {
			return nil, 0, fmt.Errorf("scan record: %w", err)
		}
		items = append(items, rec)
	}

	return items, total, nil
}

// ListRecords 根据过滤条件分页查询。
func ListRecords(filter ListFilter) ([]Record, int, error) {
	return ListRecordsWithContext(context.Background(), filter)
}

// nowMillis 返回当前时间的毫秒时间戳
func nowMillis() int64 {
	return time.Now().UnixMilli()
}
