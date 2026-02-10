package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"github.com/genwilliam/dnslog_for_go/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const apiKeyHeader = "X-API-Key"

func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil || !cfg.APIKeyRequired {
			c.Next()
			return
		}

		if c.Request.Method == http.MethodPost && (c.FullPath() == "/keys" || c.FullPath() == "/api/keys") {
			hasKeys, err := dnslog.HasAPIKeysWithContext(c.Request.Context())
			if err == nil && !hasKeys {
				c.Next()
				return
			}
		}

		key := strings.TrimSpace(c.GetHeader(apiKeyHeader))
		if key == "" {
			debugAuth(c, "", "", false, false)
			response.Error(c, http.StatusUnauthorized, response.CodeMissingAPIKey)
			c.Abort()
			return
		}

		if !isValidAPIKey(key) {
			debugAuth(c, key, dnslog.HashAPIKey(key), false, false)
			response.Error(c, http.StatusUnauthorized, response.CodeInvalidKey)
			c.Abort()
			return
		}

		hash := dnslog.HashAPIKey(key)
		apiKey, err := dnslog.GetAPIKeyByHashWithContext(c.Request.Context(), hash)
		if err != nil {
			debugAuth(c, key, hash, false, false)
			response.Error(c, http.StatusUnauthorized, response.CodeInvalidKey)
			c.Abort()
			return
		}
		if !apiKey.Enabled {
			debugAuth(c, key, hash, true, false)
			response.Error(c, http.StatusUnauthorized, response.CodeDisabledAPIKey)
			c.Abort()
			return
		}

		debugAuth(c, key, hash, true, true)
		c.Set("api_key_id", apiKey.ID)
		c.Next()

		_ = dnslog.TouchAPIKeyLastUsed(apiKey.ID, time.Now().UnixMilli())
	}
}

func isValidAPIKey(key string) bool {
	if len(key) != 64 {
		return false
	}
	for _, ch := range key {
		if (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F') {
			continue
		}
		return false
	}
	return true
}

func debugAuth(c *gin.Context, key string, hash string, hit bool, enabled bool) {
	if !gin.IsDebugging() {
		return
	}
	log.Debug("auth check",
		zap.String("path", c.FullPath()),
		zap.String("method", c.Request.Method),
		zap.Int("key_len", len(key)),
		zap.String("key_prefix", maskPrefix(key)),
		zap.String("hash_prefix", maskPrefix(hash)),
		zap.Bool("hit", hit),
		zap.Bool("enabled", enabled),
	)
}

func maskPrefix(val string) string {
	if val == "" {
		return ""
	}
	if len(val) <= 6 {
		return val
	}
	return val[:6]
}
