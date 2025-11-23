package utils

import (
	"dnslog_for_go/pkg/log"

	"gopkg.in/ini.v1"
)

func SelectPact(s string) string {
	cfg, err := ini.Load("config/dns_server.ini")
	if err != nil {
		log.Error("无法读取配置文件")
		panic("Unable to read configuration file")
	} else {
		log.Info("读取配置文件成功")
	}

	return cfg.Section("PACT").Key(s).String()
}
