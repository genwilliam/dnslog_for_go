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

	Start int64 // 起始时间戳（毫秒）
	End   int64 // 结束时间戳（毫秒）
}

var db *sql.DB

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
    INDEX idx_ts (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := conn.Exec(schema); err != nil {
		return fmt.Errorf("create table dns_records: %w", err)
	}
	return nil
}

// AddRecord 持久化一条记录。
func AddRecord(rec Record) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
INSERT INTO dns_records (domain, client_ip, protocol, qtype, timestamp, server, token)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.Domain, rec.ClientIP, rec.Protocol, rec.QType, rec.Timestamp, rec.Server, rec.Token)
	return err
}

// ListRecords 根据过滤条件分页查询。
func ListRecords(filter ListFilter) ([]Record, int, error) {
	if db == nil {
		return nil, 0, errors.New("store not initialized")
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
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

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	countSQL := "SELECT COUNT(1) FROM dns_records " + whereSQL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var total int
	if err := db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	querySQL := `
SELECT id, domain, client_ip, protocol, qtype, timestamp, server, token
FROM dns_records
` + whereSQL + `
ORDER BY timestamp DESC
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

// nowMillis 返回当前时间的毫秒时间戳
func nowMillis() int64 {
	return time.Now().UnixMilli()
}
