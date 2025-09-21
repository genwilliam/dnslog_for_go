package change_dns_server

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/ini.v1"
)

func TestReadIniFile(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Println("当前工作目录:", dir)

	cfg, err := ini.Load("dns_server.ini")
	if err != nil {
		t.Fatalf("无法读取配置文件: %v", err)
	}

	server := cfg.Section("DNS").Key("server").String()
	switch server {
	case "0":
		fmt.Println("使用 DNS：8.8.8.8")
	case "1":
		fmt.Println("使用 DNS：223.5.5.5")
	default:
		t.Errorf("未知 DNS 设置: %s", server)
	}

	err = setDNS("1")
	if err != nil {
		t.Errorf("设置 DNS 失败: %v", err)
		return
	}

	fmt.Println("DNS 设置成功")

	// 重新加载并验证
	cfg, err = ini.Load("dns_server.ini")
	if err != nil {
		t.Errorf("重新加载配置失败: %v", err)
		return
	}

	fmt.Println("修改后 DNS 值:", cfg.Section("DNS").Key("server").String())
}

func setDNS(value string) error {
	cfg, err := ini.Load("dns_server.ini")
	if err != nil {
		return err
	}

	cfg.Section("DNS").Key("server").SetValue(value)
	return cfg.SaveTo("dns_server.ini")
}
