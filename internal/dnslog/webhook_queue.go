package dnslog

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/infra"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"go.uber.org/zap"
)

const webhookQueueKey = "webhook:queue"

func MaybeEnqueueWebhook(token string, isFirst bool, domain string) error {
	cfg := config.Get()
	if cfg == nil || !cfg.WebhookEnabled {
		return nil
	}

	hook, err := GetTokenWebhook(token)
	if err != nil || !hook.Enabled {
		return nil
	}
	if hook.Mode == "FIRST_HIT" && !isFirst {
		return nil
	}
	if hook.Mode != "FIRST_HIT" && hook.Mode != "EACH_HIT" {
		return errors.New("invalid webhook mode")
	}

	payloadMap := map[string]interface{}{
		"token":     token,
		"domain":    domain,
		"hit_count": 1,
		"timestamp": time.Now().UnixMilli(),
	}
	payloadBytes, _ := json.Marshal(payloadMap)

	encSecret, err := EncryptWebhookSecret(hook.Secret)
	if err != nil {
		return err
	}
	job := WebhookJob{
		Token:       token,
		URL:         hook.URL,
		Payload:     string(payloadBytes),
		Secret:      encSecret,
		NextRetryAt: time.Now().UnixMilli(),
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
	}
	jobID, err := CreateWebhookJob(job)
	if err != nil {
		return err
	}
	return EnqueueWebhookJob(jobID)
}

func EnqueueWebhookJob(jobID int64) error {
	client := infra.GetRedis()
	if client == nil {
		return errors.New("redis not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return client.LPush(ctx, webhookQueueKey, jobID).Err()
}

func StartWebhookWorkers() {
	client := infra.GetRedis()
	if client == nil {
		return
	}
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := client.BRPop(ctx, 3*time.Second, webhookQueueKey).Result()
			cancel()
			if err != nil || len(res) < 2 {
				continue
			}
			jobID, err := parseInt64(res[1])
			if err != nil {
				continue
			}
			processWebhookJob(jobID)
		}
	}()

	go func() {
		interval := time.Duration(config.Get().WebhookRetryIntervalSeconds) * time.Second
		if interval <= 0 {
			interval = 30 * time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			nowMs := time.Now().UnixMilli()
			ids, err := ListDueWebhookJobs(nowMs, 200)
			if err != nil {
				log.Error("list due webhook jobs failed", zap.Error(err))
				continue
			}
			for _, id := range ids {
				_ = EnqueueWebhookJob(id)
			}
		}
	}()
}

func processWebhookJob(jobID int64) {
	cfg := config.Get()
	job, err := GetWebhookJob(jobID)
	if err != nil || job.Status != "PENDING" {
		return
	}
	if job.NextRetryAt > time.Now().UnixMilli() {
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	secret, err := DecryptWebhookSecret(job.Secret)
	if err != nil {
		nowMs := time.Now().UnixMilli()
		_ = UpdateWebhookJob(jobID, "FAILED", job.RetryCount, job.NextRetryAt, nowMs)
		return
	}
	req, err := http.NewRequest("POST", job.URL, bytes.NewBufferString(job.Payload))
	if err != nil {
		nowMs := time.Now().UnixMilli()
		_ = UpdateWebhookJob(jobID, "FAILED", job.RetryCount, job.NextRetryAt, nowMs)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-ID", strconv.FormatInt(jobID, 10))
	if secret != "" {
		req.Header.Set("X-Signature", signWebhook(job.Payload, secret))
	}
	resp, err := client.Do(req)
	if err == nil && resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	nowMs := time.Now().UnixMilli()
	if err == nil && resp != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_ = UpdateWebhookJob(jobID, "SUCCESS", job.RetryCount, job.NextRetryAt, nowMs)
		return
	}

	retryCount := job.RetryCount + 1
	if retryCount >= cfg.WebhookMaxRetries {
		_ = UpdateWebhookJob(jobID, "FAILED", retryCount, job.NextRetryAt, nowMs)
		return
	}

	nextRetry := nowMs + retryBackoff(retryCount)
	_ = UpdateWebhookJob(jobID, "PENDING", retryCount, nextRetry, nowMs)
}

func signWebhook(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func retryBackoff(retry int) int64 {
	switch retry {
	case 1:
		return int64(60 * 1000)
	case 2:
		return int64(5 * 60 * 1000)
	case 3:
		return int64(15 * 60 * 1000)
	default:
		return int64(60 * 60 * 1000)
	}
}

func parseInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}
