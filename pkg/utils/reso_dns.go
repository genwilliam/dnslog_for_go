package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"github.com/miekg/dns"
	"go.uber.org/zap"
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
	cfg := config.Get()
	c := &dns.Client{
		Net:     cfg.Protocol,
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
	case dns.TypeCNAME:
		for _, ans := range r.Answer {
			if cname, ok := ans.(*dns.CNAME); ok {
				results = append(results, DNSResult{
					IP:      cname.Target,
					Address: getServer(),
				})
			}
		}
	case dns.TypeMX:
		for _, ans := range r.Answer {
			if mx, ok := ans.(*dns.MX); ok {
				results = append(results, DNSResult{
					IP:      mx.Mx,
					Address: getServer(),
				})
			}
		}
	case dns.TypeTXT:
		for _, ans := range r.Answer {
			if txt, ok := ans.(*dns.TXT); ok {
				results = append(results, DNSResult{
					IP:      strings.Join(txt.Txt, "; "),
					Address: getServer(),
				})
			}
		}
	}
	return results
}

// getServer 从配置文件读取 DNS 服务器地址
func getServer() string {
	cfg := config.Get()
	return cfg.CurrentUpstream()
}
