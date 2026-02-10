package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()
	dsn := cfg.MySQLDSN
	if dsn == "" {
		fmt.Println("MYSQL_DSN is empty")
		os.Exit(1)
	}
	if os.Getenv("WEBHOOK_SECRET_KEY") == "" && cfg.WebhookSecretKey == "" {
		fmt.Println("WEBHOOK_SECRET_KEY is required")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := migrateTable(db, "token_webhooks"); err != nil {
		fmt.Println("migrate token_webhooks:", err)
		os.Exit(1)
	}
	if err := migrateTable(db, "webhook_jobs"); err != nil {
		fmt.Println("migrate webhook_jobs:", err)
		os.Exit(1)
	}

	fmt.Println("migration complete")
}

func migrateTable(db *sql.DB, table string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := fmt.Sprintf("SELECT id, secret FROM %s WHERE secret != '' AND secret NOT LIKE 'enc:%%'", table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var secret string
		if err := rows.Scan(&id, &secret); err != nil {
			return err
		}
		enc, err := dnslog.EncryptWebhookSecret(strings.TrimSpace(secret))
		if err != nil {
			return err
		}
		if enc == "" {
			continue
		}
		if err := updateSecret(db, table, id, enc); err != nil {
			return err
		}
	}
	return nil
}

func updateSecret(db *sql.DB, table string, id int64, enc string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stmt := fmt.Sprintf("UPDATE %s SET secret = ? WHERE id = ?", table)
	_, err := db.ExecContext(ctx, stmt, enc, id)
	return err
}
