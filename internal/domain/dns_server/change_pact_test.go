package dns_server_test

import (
	"fmt"
	"testing"

	"gopkg.in/ini.v1"
)

func TestIni(t *testing.T) {
	cfg, err := ini.Load("default.ini")
	if err != nil {
		panic("Unable to read configuration file")
	}

	current1 := cfg.Section("PACT").Key("udp").String()
	current2 := cfg.Section("PACT").Key("tcp").String()
	fmt.Println(current1)
	fmt.Println(current2)
}
