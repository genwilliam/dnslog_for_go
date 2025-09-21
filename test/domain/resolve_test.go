package domain

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func Test_demo(t *testing.T) {
	c := &dns.Client{
		Net:     "udp",
		Timeout: 3 * time.Second,
	}
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)

	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		t.Fatalf("Failed to exchange: %v", err)
	}
	if len(r.Answer) == 0 {
		t.Fatalf("No answer received")
	}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			t.Logf("IP Address: %s", a.A.String())
		}
	}
}
