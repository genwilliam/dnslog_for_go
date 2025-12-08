package dnslog

import (
	"github.com/genwilliam/dnslog_for_go/internal/domain/dns_server"
	"github.com/genwilliam/dnslog_for_go/pkg/log"
	"net"
	"strings"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

var (
	udpServer  *dns.Server
	tcpServer  *dns.Server
	listenAddr = ":15353"

	// TODO: 你可以把它改成从 config 读取
	rootDomain = "demo.com"
)

// StartDNSServer 启动真正的 DNS 服务器（UDP + TCP）
func StartDNSServer() {
	InitStore(1000)

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
		log.Info("DNS UDP server listening", zap.String("addr", listenAddr))
		if err := udpServer.ListenAndServe(); err != nil {
			log.Error("DNS UDP server failed", zap.Error(err))
		}
	}()

	go func() {
		log.Info("DNS TCP server listening", zap.String("addr", listenAddr))
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

	// 去掉末尾的 "."
	if strings.HasSuffix(qName, ".") {
		qName = qName[:len(qName)-1]
	}

	qType := dns.TypeToString[q.Qtype]

	// ====== 只记录 rootDomain 下的域名 ======
	if strings.HasSuffix(qName, rootDomain) {

		// 提取子域，如 abc.demo.com → abc
		subdomain := strings.TrimSuffix(qName, "."+rootDomain)
		token := "(none)"

		if subdomain != "" {
			parts := strings.Split(subdomain, ".")
			token = parts[0]
		}

		AddRecord(Record{
			Domain:    qName, // 完整域名
			ClientIP:  clientIP,
			Protocol:  proto,
			QType:     qType,
			Timestamp: nowMillis(),
			Server:    listenAddr,
			Token:     token,
		})

		log.Info("Captured DNS query",
			zap.String("domain", qName),
			zap.String("token", token),
			zap.String("qtype", qType),
			zap.String("client_ip", clientIP),
			zap.String("protocol", proto),
		)
	}

	// ====== 保留原来的上游 DNS 转发 ======
	upstreamHost := dns_server.GetDNSServer(0)
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
