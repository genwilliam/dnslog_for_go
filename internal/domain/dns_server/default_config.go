package dns_server

import (
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

func DefaultConfig() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("程序异常终止: ", zap.Any("r", r))
		}
	}()

	cfg, err := ini.Load("config/dns_server.ini")
	if err != nil {
		panic("无法读取配置文件")
	}
	cfg.Section("DNS").Key("server").SetValue("0")

	err = cfg.SaveTo("config/dns_server.ini")
	if err != nil {
		panic("默认配置恢复失败")
	} else {
		log.Info("默认配置恢复成功")
	}
}
