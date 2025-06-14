package utils

import (
	"dnslog_for_go/internal/config"
	"dnslog_for_go/internal/domain/dns_server"
	"dnslog_for_go/pkg/log"
	"fmt"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
	"strconv"
	"time"
)

// DNSQueryResult DNS 查询结果结构体
type DNSQueryResult struct {
	Domain  string      `json:"domain"`
	Results []DNSResult `json:"results"` // 用一个切片存储多个结果
}

// DNSResult 单个 DNS 记录
type DNSResult struct {
	IP      string `json:"ip"`
	Address string `json:"address"`
}

// ResolveDNS dns查询
func ResolveDNS(domainName string) DNSQueryResult {
	c := &dns.Client{
		Net:     config.GlobalPact,
		Timeout: 10 * time.Second,
	}

	// 初始化存储查询结果的切片
	var results []DNSResult

	// 定义查询类型
	queryTypes := []uint16{
		dns.TypeA,
		dns.TypeAAAA,
		dns.TypeCNAME,
		dns.TypeMX,
		dns.TypeTXT,
	}

	// 循环执行不同类型的 DNS 查询
	for _, queryType := range queryTypes {
		results = appendResults(results, domainName, queryType, c)
	}

	// 检查结果
	log.Info("查询结果", zap.Int("resultCount", len(results)))

	// 返回多个结果
	return DNSQueryResult{
		Domain:  domainName,
		Results: results, // 返回多个结果
	}
}

// appendResults 执行具体的 DNS 查询并将结果追加到 results 数组中
func appendResults(results []DNSResult, domainName string, queryType uint16, c *dns.Client) []DNSResult {
	message := new(dns.Msg)
	message.SetQuestion(dns.Fqdn(domainName), queryType) // 设置查询类型

	r, _, err := c.Exchange(message, getServer()+":53")
	if err != nil {
		log.Error(fmt.Sprintf("DNS query failed for type %d: %v", queryType, err), zap.Error(err))
		return results
	}

	// 根据查询类型处理不同的 DNS 记录
	switch queryType {
	case dns.TypeA:
		for _, ans := range r.Answer {
			if aRecord, ok := ans.(*dns.A); ok {
				results = append(results, DNSResult{
					IP:      aRecord.A.String(),
					Address: getServer(),
				})
			}
		}
	case dns.TypeAAAA:
		for _, ans := range r.Answer {
			if aaaaRecord, ok := ans.(*dns.AAAA); ok {
				results = append(results, DNSResult{
					IP:      aaaaRecord.AAAA.String(),
					Address: getServer(),
				})
			}
		}

	}
	return results
}

// getServer 从配置文件读取 DNS 服务器地址
func getServer() string {
	cfg, err := ini.Load("internal/config/dns_server.ini")
	if err != nil {
		log.Error("无法读取配置文件")
		panic("无法读取配置文件")
	}

	current := cfg.Section("DNS").Key("server").String()
	if current == "127.0.0.1" {
		return current
	}

	currentNum, err := strconv.Atoi(current)
	if err != nil {
		log.Error("配置值不是有效数字")
		panic("配置值不是有效数字")
	}
	return dns_server.DnsServer(currentNum)
}
