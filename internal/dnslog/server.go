package dnslog

import (
	"net"
	"strings"

	"github.com/genwilliam/dnslog_for_go/config"
	"github.com/genwilliam/dnslog_for_go/pkg/log"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

var (
	udpServer *dns.Server
	tcpServer *dns.Server

	listenAddr   = ":15353"
	rootDomain   = "demo.com"
	rootDomains  []string
	captureAll   bool
	activeConfig *config.Config
)

// StartDNSServer 启动 DNS 服务器（UDP + TCP）
func StartDNSServer(cfg *config.Config) {
	activeConfig = cfg
	listenAddr = cfg.DNSListenAddr
	rootDomain = strings.ToLower(strings.TrimSuffix(cfg.RootDomain, "."))
	rootDomains = make([]string, 0, len(cfg.RootDomains))
	for _, rd := range cfg.RootDomains {
		rd = strings.ToLower(strings.TrimSuffix(rd, "."))
		if rd != "" {
			rootDomains = append(rootDomains, rd)
		}
	}
	captureAll = cfg.CaptureAll
	if rootDomain == "" && len(rootDomains) > 0 {
		rootDomain = rootDomains[0]
	}
	if rootDomain == "" && !captureAll {
		rootDomain = "demo.com"
	}

	if err := InitStore(cfg.MySQLDSN); err != nil {
		log.Fatal("init store failed", zap.Error(err))
		return
	}

	dns.HandleFunc(".", handleDNSQuery)

	udpServer = &dns.Server{
		Addr: listenAddr,
		Net:  "udp",
	}
	tcpServer = &dns.Server{
		Addr: listenAddr,
		Net:  "tcp",
	}

	go func() {
		log.Info("DNS UDP server listening", zap.String("addr", listenAddr), zap.String("root_domain", rootDomain), zap.String("upstream", cfg.CurrentUpstream()))
		if err := udpServer.ListenAndServe(); err != nil {
			log.Error("DNS UDP server failed", zap.Error(err))
		}
	}()

	go func() {
		log.Info("DNS TCP server listening", zap.String("addr", listenAddr), zap.String("root_domain", rootDomain), zap.String("upstream", cfg.CurrentUpstream()))
		if err := tcpServer.ListenAndServe(); err != nil {
			log.Error("DNS TCP server failed", zap.Error(err))
		}
	}()
}

// ShutdownDNSServer 关闭 DNS 服务器
func ShutdownDNSServer() {
	if udpServer != nil {
		if err := udpServer.Shutdown(); err != nil {
			log.Error("DNS UDP shutdown failed", zap.Error(err))
		}
	}

	if tcpServer != nil {
		if err := tcpServer.Shutdown(); err != nil {
			log.Error("DNS TCP shutdown failed", zap.Error(err))
		}
	}
}

// 处理每一个 DNS 查询
func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}

	remoteAddr := w.RemoteAddr()
	clientIP := parseClientIP(remoteAddr)

	proto := "udp"
	switch remoteAddr.(type) {
	case *net.TCPAddr:
		proto = "tcp"
	}

	q := r.Question[0]
	qName := dns.Fqdn(q.Name)
	qNameLower := strings.ToLower(qName)

	// 去掉末尾的 "."
	if strings.HasSuffix(qName, ".") {
		qName = qName[:len(qName)-1]
	}
	if strings.HasSuffix(qNameLower, ".") {
		qNameLower = qNameLower[:len(qNameLower)-1]
	}

	qType := dns.TypeToString[q.Qtype]

	// ====== 是否记录该域名 ======
	matchedRoot := selectMatchedRoot(qNameLower)
	if captureAll || matchedRoot != "" {

		// 提取子域，如 abc.demo.com → abc
		subdomain := "(none)"
		if matchedRoot != "" && strings.HasSuffix(qNameLower, matchedRoot) {
			trimmed := strings.TrimSuffix(qNameLower, "."+matchedRoot)
			if trimmed != "" && trimmed != qName {
				subdomain = trimmed
			}
		}
		token := "(none)"

		if subdomain != "" && subdomain != "(none)" {
			parts := strings.Split(subdomain, ".")
			token = parts[0]
		}

		if err := AddRecord(Record{
			Domain:    qName, // 完整域名
			ClientIP:  clientIP,
			Protocol:  proto,
			QType:     qType,
			Timestamp: nowMillis(),
			Server:    listenAddr,
			Token:     token,
		}); err != nil {
			log.Error("保存 DNS 记录失败", zap.Error(err))
		}

		log.Info("Captured DNS query",
			zap.String("domain", qName),
			zap.String("token", token),
			zap.String("qtype", qType),
			zap.String("client_ip", clientIP),
			zap.String("protocol", proto),
		)
	}

	// ====== 保留原来的上游 DNS 转发 ======
	upstreamHost := activeConfig.CurrentUpstream()
	upstreamAddr := net.JoinHostPort(upstreamHost, "53")

	resp, err := dns.Exchange(r, upstreamAddr)
	if err != nil {
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}

	resp.Id = r.Id
	_ = w.WriteMsg(resp)
}

// 解析 "ip:port"
func parseClientIP(addr net.Addr) string {
	switch v := addr.(type) {
	case *net.UDPAddr:
		return v.IP.String()
	case *net.TCPAddr:
		return v.IP.String()
	default:
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return addr.String()
		}
		return host
	}
}

// selectMatchedRoot 返回与 qName 匹配的根域（若 captureAll 则返回 "" 但会被允许记录）
func selectMatchedRoot(qName string) string {
	if captureAll {
		return ""
	}
	// 优先单 rootDomain 兼容旧配置
	if rootDomain != "" && strings.HasSuffix(qName, rootDomain) {
		return rootDomain
	}
	for _, rd := range rootDomains {
		if strings.HasSuffix(qName, rd) {
			return rd
		}
	}
	return ""
}
