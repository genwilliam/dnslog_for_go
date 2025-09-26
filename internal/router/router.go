package router

import (
	"context"
	"dnslog_for_go/internal/domain"
	"dnslog_for_go/internal/domain/dns_server"
	"dnslog_for_go/pkg/log"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StartServer 启动 HTTP 服务器
func StartServer(embedFS embed.FS) {
	r := gin.Default()

	// 加载嵌入静态文件
	if err := loadStatic(r, embedFS); err != nil {
		log.Error("Failed to load static files", zap.Error(err))
		return
	}

	// 加载 HTML 模板
	if err := loadTemplates(r, embedFS); err != nil {
		log.Error("Failed to load template files", zap.Error(err))
		return
	}

	// 注册路由
	registerRoutes(r)

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 启动服务器
	go func() {
		log.Info("Server started on :8080")
		log.Info("Please visit http://localhost:8080/dnslog to access the DNS log system")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server run failed", zap.Error(err))
		}
	}()

	// 监听退出信号
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-exitChannel

	log.Info("Shutting down server gracefully...")

	// 恢复默认配置
	dns_server.DefaultConfig()

	// 优雅关闭服务器，最多等待 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", zap.Error(err))
	}

	log.Info("Server exited")
}

// loadStatic 嵌入静态文件
func loadStatic(r *gin.Engine, embedFS embed.FS) error {
	staticFiles, err := fs.Sub(embedFS, "static")
	if err != nil {
		return err
	}
	r.StaticFS("/static", http.FS(staticFiles))
	return nil
}

// loadTemplates 嵌入 HTML 模板
func loadTemplates(r *gin.Engine, embedFS embed.FS) error {
	tmplFiles, err := fs.Sub(embedFS, "templates")
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFS(tmplFiles, "*.html")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(tmpl)
	return nil
}

// registerRoutes 注册路由
func registerRoutes(r *gin.Engine) {
	r.GET("/dnslog", domain.ShowForm)
	r.POST("/submit", domain.SubmitDomain)
	r.POST("/random-domain", domain.RandomDomain)
	r.POST("/change", domain.ChangeServer)
	r.POST("/change-pact", domain.ChangePact)
	r.POST("/pause", domain.InitPause)
	r.POST("/start", domain.InitPause)
}
