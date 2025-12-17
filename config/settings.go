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

		// 1) 文件配置（config/config.yaml）
		_ = loadFromFile(appConfig)

		// 2) 环境变量覆盖
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
		RootDomain:      "demo.com",
		RootDomains:     []string{"demo.com"},
		CaptureAll:      false,
		DNSListenAddr:   ":15353",
		HTTPListenAddr:  ":8080",
		UpstreamDNS:     []string{"8.8.8.8", "223.5.5.5"},
		Protocol:        "udp",
		MySQLDSN:        "dnslog:dnslog@tcp(localhost:3306)/dnslog?parseTime=true&loc=Local&charset=utf8mb4",
		DefaultPageSize: 20,
		MaxPageSize:     100,
	}
}

func loadFromFile(cfg *Config) error {
	path := filepath.Join("config", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	type fileConfig struct {
		RootDomain     string   `yaml:"rootDomain"`
		RootDomains    []string `yaml:"rootDomains"`
		CaptureAll     bool     `yaml:"captureAll"`
		DNSListenAddr  string   `yaml:"dnsListenAddr"`
		HTTPListenAddr string   `yaml:"httpListenAddr"`
		UpstreamDNS    []string `yaml:"upstreamDNS"`
		Protocol       string   `yaml:"protocol"`
		MySQLDSN       string   `yaml:"mysqlDSN"`
		PageSize       int      `yaml:"pageSize"`
		MaxPageSize    int      `yaml:"maxPageSize"`
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
