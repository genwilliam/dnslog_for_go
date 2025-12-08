package config

import (
	"log"
	"os"
	"path/filepath"
)

// LoadDNSConfig 读取 dns_server.ini 配置文件
func LoadDNSConfig() []byte {
	// 当前工作目录
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	configPath := filepath.Join(baseDir, "config", "dns_server.ini")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", configPath, err)
	}
	return data
}
