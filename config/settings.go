package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	RootDomain     string   `yaml:"rootDomain"`     // 单个根域，兼容旧配置
	RootDomains    []string `yaml:"rootDomains"`    // 支持多个根域（可选）
	CaptureAll     bool     `yaml:"captureAll"`     // 是否记录所有域名的查询
	DNSListenAddr  string   `yaml:"dnsListenAddr"`  // DNS 监听地址，默认 :15353
	HTTPListenAddr string   `yaml:"httpListenAddr"` // HTTP 监听地址，默认 :8080
	UpstreamDNS    []string `yaml:"upstreamDNS"`    // 上游 DNS 列表
	Protocol       string   `yaml:"protocol"`       // 默认查询协议 udp/tcp
	MySQLDSN       string   `yaml:"mysqlDSN"`       // MySQL 连接串

	DefaultPageSize int `yaml:"pageSize"`
	MaxPageSize     int `yaml:"maxPageSize"`
	TokenTTLSeconds int `yaml:"tokenTTLSeconds"`

	APIKeyRequired              bool   `yaml:"apiKeyRequired"`
	BootstrapEnabled            bool   `yaml:"bootstrapEnabled"`
	BootstrapToken              string `yaml:"bootstrapToken"`
	RateLimitEnabled            bool   `yaml:"rateLimitEnabled"`
	RateLimitWindowSeconds      int    `yaml:"rateLimitWindowSeconds"`
	RateLimitMaxRequests        int    `yaml:"rateLimitMaxRequests"`
	DNSRateLimitEnabled         bool   `yaml:"dnsRateLimitEnabled"`
	DNSRateLimitWindowSeconds   int    `yaml:"dnsRateLimitWindowSeconds"`
	DNSRateLimitMaxRequests     int    `yaml:"dnsRateLimitMaxRequests"`
	RedisAddr                   string `yaml:"redisAddr"`
	RedisPassword               string `yaml:"redisPassword"`
	RedisDB                     int    `yaml:"redisDB"`
	AuditEnabled                bool   `yaml:"auditEnabled"`
	PublicConfig                bool   `yaml:"publicConfig"`
	WebhookEnabled              bool   `yaml:"webhookEnabled"`
	WebhookMaxRetries           int    `yaml:"webhookMaxRetries"`
	WebhookRetryIntervalSeconds int    `yaml:"webhookRetryIntervalSeconds"`
	WebhookSecretKey            string `yaml:"webhookSecretKey"`
	MetricsEnabled              bool   `yaml:"metricsEnabled"`
	MetricsPublic               bool   `yaml:"metricsPublic"`
	RetentionEnabled            bool   `yaml:"retentionEnabled"`
	RecordRetentionDays         int    `yaml:"recordRetentionDays"`
	RetentionIntervalSeconds    int    `yaml:"retentionIntervalSeconds"`
	RetentionBatchSize          int    `yaml:"retentionBatchSize"`

	currentUpstream int
	mu              sync.RWMutex
}

var (
	appConfig *Config
	once      sync.Once
)

// Load 初始化配置，只会执行一次。
func Load() *Config {
	once.Do(func() {
		appConfig = defaultConfig()

		// 文件配置（config/config.yaml）
		_ = loadFromFile(appConfig)

		// 环境变量覆盖
		applyEnvOverrides(appConfig)

		if appConfig.Protocol != "udp" && appConfig.Protocol != "tcp" {
			appConfig.Protocol = "udp"
		}

		// 兼容 RootDomain 单值
		appConfig.RootDomains = mergeRootDomains(appConfig.RootDomain, appConfig.RootDomains)

		// 兼容旧逻辑
	})
	return appConfig
}

// Get 返回已加载的配置，若未加载则 panic。
func Get() *Config {
	if appConfig == nil {
		return Load()
	}
	return appConfig
}

// CurrentUpstream 返回当前使用的上游 DNS。
func (c *Config) CurrentUpstream() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.UpstreamDNS) == 0 {
		return "8.8.8.8"
	}
	if c.currentUpstream < 0 || c.currentUpstream >= len(c.UpstreamDNS) {
		return c.UpstreamDNS[0]
	}
	return c.UpstreamDNS[c.currentUpstream]
}

// GetProtocol 返回当前的 DNS 协议
func (c *Config) GetProtocol() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Protocol
}

// SetProtocol 更新 DNS 协议
func (c *Config) SetProtocol(p string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Protocol = p
}

// SetUpstreamIndex 更新上游 DNS 下标。
func (c *Config) SetUpstreamIndex(idx int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if idx >= 0 && idx < len(c.UpstreamDNS) {
		c.currentUpstream = idx
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func splitAndTrim(val string) []string {
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func mustInt(val string, def int) int {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return def
}

func defaultConfig() *Config {
	return &Config{
		RootDomain:                  "demo.com",
		RootDomains:                 []string{"demo.com"},
		CaptureAll:                  false,
		DNSListenAddr:               ":15353",
		HTTPListenAddr:              ":8080",
		UpstreamDNS:                 []string{"8.8.8.8", "223.5.5.5"},
		Protocol:                    "udp",
		MySQLDSN:                    "dnslog:dnslog@tcp(localhost:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4",
		DefaultPageSize:             20,
		MaxPageSize:                 100,
		TokenTTLSeconds:             3600,
		APIKeyRequired:              true,
		BootstrapEnabled:            false,
		BootstrapToken:              "",
		RateLimitEnabled:            true,
		RateLimitWindowSeconds:      60,
		RateLimitMaxRequests:        60,
		DNSRateLimitEnabled:         true,
		DNSRateLimitWindowSeconds:   60,
		DNSRateLimitMaxRequests:     1000,
		RedisAddr:                   "127.0.0.1:6379",
		RedisPassword:               "",
		RedisDB:                     0,
		AuditEnabled:                true,
		PublicConfig:                false,
		WebhookEnabled:              true,
		WebhookMaxRetries:           4,
		WebhookRetryIntervalSeconds: 30,
		WebhookSecretKey:            "",
		MetricsEnabled:              true,
		MetricsPublic:               false,
		RetentionEnabled:            true,
		RecordRetentionDays:         30,
		RetentionIntervalSeconds:    3600,
		RetentionBatchSize:          1000,
	}
}

func loadFromFile(cfg *Config) error {
	path := filepath.Join("config", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	type fileConfig struct {
		RootDomain                  string   `yaml:"rootDomain"`
		RootDomains                 []string `yaml:"rootDomains"`
		CaptureAll                  bool     `yaml:"captureAll"`
		DNSListenAddr               string   `yaml:"dnsListenAddr"`
		HTTPListenAddr              string   `yaml:"httpListenAddr"`
		UpstreamDNS                 []string `yaml:"upstreamDNS"`
		Protocol                    string   `yaml:"protocol"`
		MySQLDSN                    string   `yaml:"mysqlDSN"`
		PageSize                    int      `yaml:"pageSize"`
		MaxPageSize                 int      `yaml:"maxPageSize"`
		TokenTTLSeconds             int      `yaml:"tokenTTLSeconds"`
		APIKeyRequired              *bool    `yaml:"apiKeyRequired"`
		BootstrapEnabled            *bool    `yaml:"bootstrapEnabled"`
		BootstrapToken              string   `yaml:"bootstrapToken"`
		RateLimitEnabled            *bool    `yaml:"rateLimitEnabled"`
		RateLimitWindowSeconds      int      `yaml:"rateLimitWindowSeconds"`
		RateLimitMaxRequests        int      `yaml:"rateLimitMaxRequests"`
		DNSRateLimitEnabled         *bool    `yaml:"dnsRateLimitEnabled"`
		DNSRateLimitWindowSeconds   int      `yaml:"dnsRateLimitWindowSeconds"`
		DNSRateLimitMaxRequests     int      `yaml:"dnsRateLimitMaxRequests"`
		RedisAddr                   string   `yaml:"redisAddr"`
		RedisPassword               string   `yaml:"redisPassword"`
		RedisDB                     int      `yaml:"redisDB"`
		AuditEnabled                *bool    `yaml:"auditEnabled"`
		PublicConfig                *bool    `yaml:"publicConfig"`
		WebhookEnabled              *bool    `yaml:"webhookEnabled"`
		WebhookMaxRetries           int      `yaml:"webhookMaxRetries"`
		WebhookRetryIntervalSeconds int      `yaml:"webhookRetryIntervalSeconds"`
		WebhookSecretKey            string   `yaml:"webhookSecretKey"`
		MetricsEnabled              *bool    `yaml:"metricsEnabled"`
		MetricsPublic               *bool    `yaml:"metricsPublic"`
		RetentionEnabled            *bool    `yaml:"retentionEnabled"`
		RecordRetentionDays         int      `yaml:"recordRetentionDays"`
		RetentionIntervalSeconds    int      `yaml:"retentionIntervalSeconds"`
		RetentionBatchSize          int      `yaml:"retentionBatchSize"`
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return err
	}

	if fc.RootDomain != "" {
		cfg.RootDomain = fc.RootDomain
	}
	if len(fc.RootDomains) > 0 {
		cfg.RootDomains = fc.RootDomains
	}
	cfg.CaptureAll = fc.CaptureAll
	if fc.DNSListenAddr != "" {
		cfg.DNSListenAddr = fc.DNSListenAddr
	}
	if fc.HTTPListenAddr != "" {
		cfg.HTTPListenAddr = fc.HTTPListenAddr
	}
	if len(fc.UpstreamDNS) > 0 {
		cfg.UpstreamDNS = fc.UpstreamDNS
	}
	if fc.Protocol != "" {
		cfg.Protocol = strings.ToLower(fc.Protocol)
	}
	if fc.MySQLDSN != "" {
		cfg.MySQLDSN = fc.MySQLDSN
	}
	if fc.PageSize > 0 {
		cfg.DefaultPageSize = fc.PageSize
	}
	if fc.MaxPageSize > 0 {
		cfg.MaxPageSize = fc.MaxPageSize
	}
	if fc.TokenTTLSeconds > 0 {
		cfg.TokenTTLSeconds = fc.TokenTTLSeconds
	}
	if fc.APIKeyRequired != nil {
		cfg.APIKeyRequired = *fc.APIKeyRequired
	}
	if fc.BootstrapEnabled != nil {
		cfg.BootstrapEnabled = *fc.BootstrapEnabled
	}
	if fc.BootstrapToken != "" {
		cfg.BootstrapToken = fc.BootstrapToken
	}
	if fc.RateLimitEnabled != nil {
		cfg.RateLimitEnabled = *fc.RateLimitEnabled
	}
	if fc.RateLimitWindowSeconds > 0 {
		cfg.RateLimitWindowSeconds = fc.RateLimitWindowSeconds
	}
	if fc.RateLimitMaxRequests > 0 {
		cfg.RateLimitMaxRequests = fc.RateLimitMaxRequests
	}
	if fc.DNSRateLimitEnabled != nil {
		cfg.DNSRateLimitEnabled = *fc.DNSRateLimitEnabled
	}
	if fc.DNSRateLimitWindowSeconds > 0 {
		cfg.DNSRateLimitWindowSeconds = fc.DNSRateLimitWindowSeconds
	}
	if fc.DNSRateLimitMaxRequests > 0 {
		cfg.DNSRateLimitMaxRequests = fc.DNSRateLimitMaxRequests
	}
	if fc.RedisAddr != "" {
		cfg.RedisAddr = fc.RedisAddr
	}
	if fc.RedisPassword != "" {
		cfg.RedisPassword = fc.RedisPassword
	}
	if fc.RedisDB > 0 {
		cfg.RedisDB = fc.RedisDB
	}
	if fc.AuditEnabled != nil {
		cfg.AuditEnabled = *fc.AuditEnabled
	}
	if fc.PublicConfig != nil {
		cfg.PublicConfig = *fc.PublicConfig
	}
	if fc.WebhookEnabled != nil {
		cfg.WebhookEnabled = *fc.WebhookEnabled
	}
	if fc.WebhookMaxRetries > 0 {
		cfg.WebhookMaxRetries = fc.WebhookMaxRetries
	}
	if fc.WebhookRetryIntervalSeconds > 0 {
		cfg.WebhookRetryIntervalSeconds = fc.WebhookRetryIntervalSeconds
	}
	if fc.WebhookSecretKey != "" {
		cfg.WebhookSecretKey = fc.WebhookSecretKey
	}
	if fc.MetricsEnabled != nil {
		cfg.MetricsEnabled = *fc.MetricsEnabled
	}
	if fc.MetricsPublic != nil {
		cfg.MetricsPublic = *fc.MetricsPublic
	}
	if fc.RetentionEnabled != nil {
		cfg.RetentionEnabled = *fc.RetentionEnabled
	}
	if fc.RecordRetentionDays > 0 {
		cfg.RecordRetentionDays = fc.RecordRetentionDays
	}
	if fc.RetentionIntervalSeconds > 0 {
		cfg.RetentionIntervalSeconds = fc.RetentionIntervalSeconds
	}
	if fc.RetentionBatchSize > 0 {
		cfg.RetentionBatchSize = fc.RetentionBatchSize
	}
	return nil
}

func applyEnvOverrides(cfg *Config) {
	if v := getEnv("ROOT_DOMAIN", ""); v != "" {
		cfg.RootDomain = v
	}
	if v := getEnv("ROOT_DOMAINS", ""); v != "" {
		cfg.RootDomains = splitAndTrim(v)
	}
	if v := getEnv("CAPTURE_ALL", ""); v != "" {
		cfg.CaptureAll = strings.ToLower(v) == "true"
	}
	if v := getEnv("DNS_LISTEN_ADDR", ""); v != "" {
		cfg.DNSListenAddr = v
	}
	if v := getEnv("HTTP_LISTEN_ADDR", ""); v != "" {
		cfg.HTTPListenAddr = v
	}
	if v := getEnv("UPSTREAM_DNS", ""); v != "" {
		cfg.UpstreamDNS = splitAndTrim(v)
	}
	if v := getEnv("DNS_PROTOCOL", ""); v != "" {
		cfg.Protocol = strings.ToLower(v)
	}
	if v := getEnv("MYSQL_DSN", ""); v != "" {
		cfg.MySQLDSN = v
	}
	if v := getEnv("PAGE_SIZE", ""); v != "" {
		cfg.DefaultPageSize = mustInt(v, cfg.DefaultPageSize)
	}
	if v := getEnv("MAX_PAGE_SIZE", ""); v != "" {
		cfg.MaxPageSize = mustInt(v, cfg.MaxPageSize)
	}
	if v := getEnv("TOKEN_TTL_SECONDS", ""); v != "" {
		cfg.TokenTTLSeconds = mustInt(v, cfg.TokenTTLSeconds)
	}
	if v := getEnv("API_KEY_REQUIRED", ""); v != "" {
		cfg.APIKeyRequired = strings.ToLower(v) == "true"
	}
	if v := getEnv("BOOTSTRAP_ENABLED", ""); v != "" {
		cfg.BootstrapEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("BOOTSTRAP_TOKEN", ""); v != "" {
		cfg.BootstrapToken = v
	}
	if v := getEnv("RATE_LIMIT_ENABLED", ""); v != "" {
		cfg.RateLimitEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("RATE_LIMIT_WINDOW_SECONDS", ""); v != "" {
		cfg.RateLimitWindowSeconds = mustInt(v, cfg.RateLimitWindowSeconds)
	}
	if v := getEnv("RATE_LIMIT_MAX_REQUESTS", ""); v != "" {
		cfg.RateLimitMaxRequests = mustInt(v, cfg.RateLimitMaxRequests)
	}
	if v := getEnv("DNS_RATE_LIMIT_ENABLED", ""); v != "" {
		cfg.DNSRateLimitEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("DNS_RATE_LIMIT_WINDOW_SECONDS", ""); v != "" {
		cfg.DNSRateLimitWindowSeconds = mustInt(v, cfg.DNSRateLimitWindowSeconds)
	}
	if v := getEnv("DNS_RATE_LIMIT_MAX_REQUESTS", ""); v != "" {
		cfg.DNSRateLimitMaxRequests = mustInt(v, cfg.DNSRateLimitMaxRequests)
	}
	if v := getEnv("REDIS_ADDR", ""); v != "" {
		cfg.RedisAddr = v
	}
	if v := getEnv("REDIS_PASSWORD", ""); v != "" {
		cfg.RedisPassword = v
	}
	if v := getEnv("REDIS_DB", ""); v != "" {
		cfg.RedisDB = mustInt(v, cfg.RedisDB)
	}
	if v := getEnv("AUDIT_ENABLED", ""); v != "" {
		cfg.AuditEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("PUBLIC_CONFIG", ""); v != "" {
		cfg.PublicConfig = strings.ToLower(v) == "true"
	}
	if v := getEnv("WEBHOOK_ENABLED", ""); v != "" {
		cfg.WebhookEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("WEBHOOK_MAX_RETRIES", ""); v != "" {
		cfg.WebhookMaxRetries = mustInt(v, cfg.WebhookMaxRetries)
	}
	if v := getEnv("WEBHOOK_RETRY_INTERVAL_SECONDS", ""); v != "" {
		cfg.WebhookRetryIntervalSeconds = mustInt(v, cfg.WebhookRetryIntervalSeconds)
	}
	if v := getEnv("WEBHOOK_SECRET_KEY", ""); v != "" {
		cfg.WebhookSecretKey = v
	}
	if v := getEnv("METRICS_ENABLED", ""); v != "" {
		cfg.MetricsEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("METRICS_PUBLIC", ""); v != "" {
		cfg.MetricsPublic = strings.ToLower(v) == "true"
	}
	if v := getEnv("RETENTION_ENABLED", ""); v != "" {
		cfg.RetentionEnabled = strings.ToLower(v) == "true"
	}
	if v := getEnv("RECORD_RETENTION_DAYS", ""); v != "" {
		cfg.RecordRetentionDays = mustInt(v, cfg.RecordRetentionDays)
	}
	if v := getEnv("RETENTION_INTERVAL_SECONDS", ""); v != "" {
		cfg.RetentionIntervalSeconds = mustInt(v, cfg.RetentionIntervalSeconds)
	}
	if v := getEnv("RETENTION_BATCH_SIZE", ""); v != "" {
		cfg.RetentionBatchSize = mustInt(v, cfg.RetentionBatchSize)
	}
}

// Validate basic fields (currently minimal)
func (c *Config) Validate() error {
	if !c.CaptureAll && len(c.RootDomains) == 0 && c.RootDomain == "" {
		return errors.New("root domain is empty (or set CAPTURE_ALL=true)")
	}
	if len(c.UpstreamDNS) == 0 {
		return errors.New("upstream DNS is empty")
	}
	return nil
}

// mergeRootDomains 合并单值 RootDomain 与 RootDomains 切片，去重
func mergeRootDomains(single string, list []string) []string {
	all := make([]string, 0)
	if single != "" {
		all = append(all, single)
	}
	all = append(all, list...)
	seen := make(map[string]struct{})
	out := make([]string, 0, len(all))
	for _, v := range all {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
