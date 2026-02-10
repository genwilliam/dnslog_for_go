package router

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/dnslog"
	"github.com/genwilliam/dnslog_for_go/internal/domain"
	"github.com/genwilliam/dnslog_for_go/internal/infra"
	"github.com/genwilliam/dnslog_for_go/internal/middleware"
	"github.com/genwilliam/dnslog_for_go/internal/metrics"
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// StartServer 启动 HTTP 服务器 + DNSLog 服务器
func StartServer(cfg *config.Config) {
	r := gin.Default()

	// 添加跨域中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理浏览器的预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	if cfg.RateLimitEnabled || cfg.DNSRateLimitEnabled || cfg.AuditEnabled || cfg.WebhookEnabled {
		if _, err := infra.InitRedis(cfg); err != nil {
			log.Fatal("init redis failed", zap.Error(err))
			return
		}
	}
	if cfg.MetricsEnabled {
		metrics.Init()
	}
	if cfg.AuditEnabled {
		dnslog.StartAuditWorker()
	}
	if cfg.WebhookEnabled {
		dnslog.StartWebhookWorkers()
	}

	// 注册路由
	registerRoutes(r, cfg)

	// 启动 DNSLog 服务器（监听 :5353，捕获真实 DNS 请求）
	dnslog.StartDNSServer(cfg)

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:              cfg.HTTPListenAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 启动 HTTP 服务器
	go func() {
		log.Info("Server started", zap.String("addr", cfg.HTTPListenAddr))
		log.Info("Please visit /dnslog to access the DNS log system")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server run failed", zap.Error(err))
		}
	}()

	// 监听退出信号
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-exitChannel

	log.Info("Shutting down server gracefully...")

	// 关闭 DNSLog 服务器
	dnslog.ShutdownDNSServer()

	// 关闭 HTTP 服务器，最多等待 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", zap.Error(err))
	}

	log.Info("Server exited")
}

// registerRoutes 注册路由
func registerRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/dnslog", domain.ShowForm)
	registerAPIRoutes(r, cfg, "")
	registerAPIRoutes(r, cfg, "/api")
}

func registerAPIRoutes(r *gin.Engine, cfg *config.Config, prefix string) {
	base := r.Group(prefix)
	secured := base.Group("/")
	secured.Use(
		middleware.TraceID(),
		middleware.Audit(cfg),
		middleware.IPBlacklist(cfg),
		middleware.APIKeyAuth(cfg),
		middleware.RateLimit(cfg),
		middleware.Metrics(),
	)

	secured.POST("/submit", domain.SubmitDomain) // legacy: DNSLog 记录查询接口（观测模式）
	secured.GET("/random-domain", domain.RandomDomain)
	secured.POST("/tokens", domain.RandomDomain)
	secured.POST("/change", domain.ChangeServer)
	secured.POST("/change-pact", domain.ChangePact)
	secured.POST("/pause", domain.InitPause)
	secured.POST("/start", domain.InitPause)

	secured.GET("/records", dnslog.ListRecordsHandler)
	secured.GET("/tokens", dnslog.ListTokensHandler)
	secured.GET("/tokens/:token", dnslog.GetTokenStatusHandler)
	secured.GET("/tokens/:token/records", dnslog.GetTokenRecordsHandler)
	secured.POST("/tokens/:token/webhook", dnslog.SetTokenWebhookHandler)
	secured.GET("/tokens/:token/webhook", dnslog.GetTokenWebhookHandler)
	secured.POST("/tokens/:token/webhook/disable", dnslog.DisableTokenWebhookHandler)
	secured.DELETE("/tokens/:token/webhook", dnslog.DisableTokenWebhookHandler)
	secured.POST("/keys", dnslog.CreateAPIKeyHandler)
	if cfg != nil && cfg.BootstrapEnabled {
		base.POST("/keys/bootstrap", dnslog.CreateAPIKeyWithBootstrapHandler)
	}
	secured.GET("/keys", dnslog.ListAPIKeysHandler)
	secured.POST("/keys/:id/disable", dnslog.DisableAPIKeyHandler)
	secured.DELETE("/keys/:id", dnslog.DisableAPIKeyHandler)
	secured.POST("/blacklist", dnslog.AddBlacklistHandler)
	secured.GET("/blacklist", dnslog.ListBlacklistHandler)
	secured.POST("/blacklist/:id/disable", dnslog.DisableBlacklistHandler)
	secured.DELETE("/blacklist/:id", dnslog.DisableBlacklistHandler)

	if cfg.PublicConfig {
		base.GET("/config", ConfigHandler)
	} else {
		secured.GET("/config", ConfigHandler)
	}

	if cfg.MetricsEnabled {
		if cfg.MetricsPublic {
			base.GET("/metrics", gin.WrapH(promhttp.Handler()))
		} else {
			secured.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}
	}
}
