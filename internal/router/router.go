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
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StartServer 启动 HTTP 服务器 + DNSLog 服务器
func StartServer(cfg *config.Config) {
	r := gin.Default()

	// 添加跨域中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理浏览器的预检请求（OPTIONS）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 注册路由
	registerRoutes(r)

	// 启动 DNSLog 服务器（监听 :5353，捕获真实 DNS 请求）
	dnslog.StartDNSServer(cfg)

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    cfg.HTTPListenAddr,
		Handler: r,
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
func registerRoutes(r *gin.Engine) {
	r.GET("/dnslog", domain.ShowForm)
	r.POST("/submit", domain.SubmitDomain) // HTTP DNS 查询接口
	r.GET("/random-domain", domain.RandomDomain)
	r.POST("/change", domain.ChangeServer)
	r.POST("/change-pact", domain.ChangePact)
	r.POST("/pause", domain.InitPause)
	r.POST("/start", domain.InitPause)

	r.GET("/records", dnslog.ListRecordsHandler)
	r.GET("/config", ConfigHandler)
}
