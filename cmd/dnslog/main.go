package main

import (
	"dnslog_for_go/internal/router"
	"dnslog_for_go/pkg/log"
	"dnslog_for_go/web"
)

func main() {
	log.InitZapLogger() // 初始化日志
	router.StartServer(web.EmbedFiles)
}
