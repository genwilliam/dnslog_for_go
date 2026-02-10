package dnslog

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TokenWebhook struct {
	ID        int64
	Token     string
	URL       string
	Secret    string
	Mode      string
	Enabled   bool
	CreatedAt int64
}

type WebhookJob struct {
	ID         int64
	Token      string
	URL        string
	Payload    string
	Secret     string
	Status     string
	RetryCount int
	NextRetryAt int64
	CreatedAt  int64
	UpdatedAt  int64
}

var ErrWebhookNotFound = errors.New("webhook_not_found")

func UpsertTokenWebhookWithContext(ctx context.Context, token, url, secret, mode string, nowMs int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	encSecret, err := EncryptWebhookSecret(secret)
	if err != nil {
		return err
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, `
INSERT INTO token_webhooks (token, webhook_url, secret, mode, enabled, created_at)
VALUES (?, ?, ?, ?, 1, ?)
ON DUPLICATE KEY UPDATE webhook_url = VALUES(webhook_url), secret = VALUES(secret), mode = VALUES(mode), enabled = 1
`, token, url, encSecret, mode, nowMs)
	return err
}

func UpsertTokenWebhook(token, url, secret, mode string, nowMs int64) error {
	return UpsertTokenWebhookWithContext(context.Background(), token, url, secret, mode, nowMs)
}

func GetTokenWebhookWithContext(ctx context.Context, token string) (TokenWebhook, error) {
	if db == nil {
		return TokenWebhook{}, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var w TokenWebhook
	var enabled int
	err := db.QueryRowContext(ctx, `
SELECT id, token, webhook_url, secret, mode, enabled, created_at
FROM token_webhooks
WHERE token = ?
`, token).Scan(&w.ID, &w.Token, &w.URL, &w.Secret, &w.Mode, &enabled, &w.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return TokenWebhook{}, ErrWebhookNotFound
	}
	w.Enabled = enabled == 1
	if w.Secret != "" {
		plain, err := DecryptWebhookSecret(w.Secret)
		if err != nil {
			return TokenWebhook{}, err
		}
		w.Secret = plain
	}
	return w, err
}

func GetTokenWebhook(token string) (TokenWebhook, error) {
	return GetTokenWebhookWithContext(context.Background(), token)
}

func DisableTokenWebhookWithContext(ctx context.Context, token string) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `UPDATE token_webhooks SET enabled = 0 WHERE token = ?`, token)
	return err
}

func DisableTokenWebhook(token string) error {
	return DisableTokenWebhookWithContext(context.Background(), token)
}

func CreateWebhookJobWithContext(ctx context.Context, job WebhookJob) (int64, error) {
	if db == nil {
		return 0, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, `
INSERT INTO webhook_jobs (token, url, payload, secret, status, retry_count, next_retry_at, created_at, updated_at)
VALUES (?, ?, ?, ?, 'PENDING', 0, ?, ?, ?)
`, job.Token, job.URL, job.Payload, job.Secret, job.NextRetryAt, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func CreateWebhookJob(job WebhookJob) (int64, error) {
	return CreateWebhookJobWithContext(context.Background(), job)
}

func GetWebhookJobWithContext(ctx context.Context, id int64) (WebhookJob, error) {
	if db == nil {
		return WebhookJob{}, errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var job WebhookJob
	err := db.QueryRowContext(ctx, `
SELECT id, token, url, payload, secret, status, retry_count, next_retry_at, created_at, updated_at
FROM webhook_jobs
WHERE id = ?
`, id).Scan(&job.ID, &job.Token, &job.URL, &job.Payload, &job.Secret, &job.Status, &job.RetryCount, &job.NextRetryAt, &job.CreatedAt, &job.UpdatedAt)
	return job, err
}

func GetWebhookJob(id int64) (WebhookJob, error) {
	return GetWebhookJobWithContext(context.Background(), id)
}

func UpdateWebhookJobWithContext(ctx context.Context, id int64, status string, retryCount int, nextRetryAt, updatedAt int64) error {
	if db == nil {
		return errors.New("store not initialized")
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, `
UPDATE webhook_jobs
SET status = ?, retry_count = ?, next_retry_at = ?, updated_at = ?
WHERE id = ?
`, status, retryCount, nextRetryAt, updatedAt, id)
	return err
}

func UpdateWebhookJob(id int64, status string, retryCount int, nextRetryAt, updatedAt int64) error {
	return UpdateWebhookJobWithContext(context.Background(), id, status, retryCount, nextRetryAt, updatedAt)
}

func ListDueWebhookJobsWithContext(ctx context.Context, nowMs int64, limit int) ([]int64, error) {
	if db == nil {
		return nil, errors.New("store not initialized")
	}
	if limit <= 0 {
		limit = 200
	}
	ctx = ensureContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, `
SELECT id FROM webhook_jobs
WHERE status = 'PENDING' AND next_retry_at <= ?
ORDER BY next_retry_at ASC
LIMIT ?
`, nowMs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ListDueWebhookJobs(nowMs int64, limit int) ([]int64, error) {
	return ListDueWebhookJobsWithContext(context.Background(), nowMs, limit)
}
