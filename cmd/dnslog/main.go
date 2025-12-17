package main

import (
	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/internal/router"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
)

func main() {
	cfg := config.Load() // 加载配置
	log.InitZapLogger()  // 初始化日志
	//defer log.Sync()
	router.StartServer(cfg)

}
