package dns_server

import (
	"log"
	"strconv"

	"github.com/genwilliam/dnslog_for_go/config"

	"gopkg.in/ini.v1"
)

// GetDNSServer 根据索引返回 DNS 服务器
func GetDNSServer(num int) string {
	// 默认 DNS 列表
	servers := []string{
		"8.8.8.8",   // Google Public DNS
		"223.5.5.5", // 阿里公共 DNS
		"127.0.0.1", // 本地 DNS
	}

	// 尝试读取配置文件
	data := config.LoadDNSConfig() // 通过相对路径读取 config/dns_server.ini
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, data)
	if err != nil {
		log.Println("无法读取配置文件，使用默认 DNS:", err)
		return servers[num]
	}

	// 从配置文件读取 server 值
	serverStr := cfg.Section("DNS").Key("server").String()
	idx, err := strconv.Atoi(serverStr)
	if err != nil || idx < 0 || idx >= len(servers) {
		log.Println("配置文件 server 值无效，使用默认 DNS:", serverStr)
		return servers[num]
	}

	return servers[idx]
}
