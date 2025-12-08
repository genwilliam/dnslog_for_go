package router_test

import (
	"github.com/genwilliam/dnslog_for_go/internal/domain"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	// 配置静态文件
	r.Static("/static", "./static")

	// 加载 HTML 模板
	r.LoadHTMLGlob("templates/*")

	// 路由设置
	r.GET("/", domain.ShowForm)
	r.POST("/submit", domain.SubmitDomain)

	// 启动服务器
	r.Run(":8080")
}
