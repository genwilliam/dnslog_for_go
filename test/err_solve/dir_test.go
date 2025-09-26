package err_solve

import (
	"dnslog_for_go/pkg/log"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

func Test(t *testing.T) {
	ChangeServer(1)
}

func ChangeServer(num byte) {
	defer func() {
		if r := recover(); r != nil { // panic传什么值，recover就返回什么值
			log.Error("程序异常终止: ", zap.Any("r", r))
		}
	}()

	dir, _ := os.Getwd()
	fmt.Println(dir)

	exePath, _ := os.Executable()
	basePath := filepath.Dir(exePath)
	cfgPath := filepath.Join(basePath, "internal", "config", "dns_server.ini")
	fmt.Println(cfgPath)
	cfg, err := ini.Load(cfgPath)
	if err != nil {
		fmt.Println("无法读取配置文件")
		panic("Unable to read configuration file")
	}
	fmt.Println("读取配置文件成功")

	current := cfg.Section("DNS").Key("server").String()
	currentNum, err := strconv.Atoi(current)
	if err != nil {
		fmt.Println("配置值不是有效数字")
		panic("Configuration values are not valid numbers")
	}

	if int(num) == currentNum {
		fmt.Println("DNS 设置已是当前值，无需修改")
		return
	}

}
