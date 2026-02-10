package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	APIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dnslog_api_requests_total",
			Help: "Total HTTP API requests",
		},
		[]string{"path", "method", "status"},
	)
	DNSQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dnslog_dns_queries_total",
			Help: "Total DNS queries received",
		},
		[]string{"protocol"},
	)
	TokenHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dnslog_token_hits_total",
			Help: "Total token hits recorded",
		},
	)
)

func Init() {
	prometheus.MustRegister(APIRequestsTotal, DNSQueriesTotal, TokenHitsTotal)
}
