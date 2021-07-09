package fns

import (
	"strings"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type mockRepo struct {
	rrs []record
}

type record struct {
	RR
	Type uint16
}

var repo = mockRepo{
	[]record{
		{RR{"@", "1.1.1.1", "example.com", 600}, dns.TypeA},
		{RR{"foo", "1.1.1.1", "example.com", 600}, dns.TypeA},
		{RR{"foo", "1.1.1.2", "example.com", 600}, dns.TypeA},
		{RR{"*", "1.1.1.3", "example.com", 3600}, dns.TypeA},
		{RR{"www", "bar.example.com.", "example.com", 3600}, dns.TypeCNAME},
		{RR{"foo.bar", "cname.example.com.", "example.com", 3600}, dns.TypeCNAME},
		{RR{"@", "20 mail.vip.example.com.", "example.com", 600}, dns.TypeMX},
		{RR{"@", "10 mail.example.com.", "example.com", 600}, dns.TypeMX},
		{RR{"ipv6", "2001:4860:4860::8888", "example.com", 600}, dns.TypeAAAA},
		{RR{"@", "version=1.0", "example.com", 86400}, dns.TypeTXT},
		{RR{"@", "10 20 5060 bar.example.com.", "example.com", 600}, dns.TypeSRV},
		{RR{"@", "0 issue \"letsencrypt.org\"", "example.com", 600}, dns.TypeCAA},
		{RR{"@", "ns1.example.com.", "example.com", 1800}, dns.TypeNS},
		{RR{"@", "ns2.example.com.", "example.com", 1800}, dns.TypeNS},
	},
}

func TestResolve(t *testing.T) {
	var (
		tests = []test.Case{
			{
				Qname: "example.com.", Qtype: dns.TypeA,
				Answer: []dns.RR{test.A("example.com. 600 IN A 1.1.1.1")},
			},
			{
				Qname: "foo.example.com.", Qtype: dns.TypeA, Answer: []dns.RR{
					test.A("foo.example.com. 600 IN A 1.1.1.1"),
					test.A("foo.example.com. 600 IN A 1.1.1.2"),
				},
			},
			{
				Qname: "more.example.com.", Qtype: dns.TypeA,
				Answer: []dns.RR{test.A("more.example.com. 3600 IN A 1.1.1.3")},
			},
			{
				Qname: "www.example.com.", Qtype: dns.TypeA,
				Answer: []dns.RR{test.CNAME("www.example.com. 3600 IN CNAME bar.example.com.")},
			},
			{
				Qname: "foo.bar.example.com.", Qtype: dns.TypeA,
				Answer: []dns.RR{test.CNAME("foo.bar.example.com. 3600 IN CNAME cname.example.com.")},
			},
			{
				Qname: "foo.bar.example.com.", Qtype: dns.TypeCNAME,
				Answer: []dns.RR{test.CNAME("foo.bar.example.com. 3600 IN CNAME cname.example.com.")},
			},
			{
				Qname: "example.com.", Qtype: dns.TypeMX, Answer: []dns.RR{
					test.MX("example.com. 600 IN MX 10 mail.example.com."),
					test.MX("example.com. 600 IN MX 20 mail.vip.example.com."),
				},
			},
			{
				Qname: "ipv6.example.com.", Qtype: dns.TypeAAAA,
				Answer: []dns.RR{test.AAAA("ipv6.example.com. 600 IN AAAA 2001:4860:4860::8888")},
			},
			{
				Qname: "example.com.", Qtype: dns.TypeTXT,
				Answer: []dns.RR{test.TXT("example.com. 86400 IN TXT version=1.0")},
			},
			{
				Qname: "example.com.", Qtype: dns.TypeSRV,
				Answer: []dns.RR{test.SRV("example.com. 600 IN SRV 10 20 5060 bar.example.com.")},
			},
			{
				Qname: "example.net.", Qtype: dns.TypeAAAA,
				Answer: []dns.RR{},
			},
			{
				Qname: "example.com.", Qtype: dns.TypeNS,
				Answer: []dns.RR{
					test.NS("example.com. 1800 IN NS ns1.example.com."),
					test.NS("example.com. 1800 IN NS ns2.example.com."),
				},
			},
		}

		f = FNS{RRsRepo: repo}
	)

	for _, tc := range tests {
		r := tc.Msg()
		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		state := request.Request{W: rec, Req: r}
		a := new(dns.Msg)
		a.SetReply(r)
		a.Compress = true
		a.Authoritative = true

		a.Answer = f.Resolve(state.QName(), state.QType(), state.QClass())
		state.SizeAndDo(a)
		if err := rec.WriteMsg(a); err != nil {
			t.Error(err)
		}

		resp := rec.Msg
		if err := test.SortAndCheck(resp, tc); err != nil {
			t.Error(err)
		}
	}
}

func (repo mockRepo) GetRRs(cond Condition) (rrs []RR) {
	var (
		rrType       uint16
		domain       string
		name         string
		wildCardName string
	)

	rrs = make([]RR, 0)
	if qname, ok := cond["name"].(string); ok {
		q := resolveName(qname)

		name = q.GetName()
		wildCardName = q.GetWildCardName()
		domain = q.GetDomain()
	}

	if qType, ok := cond["type"].(uint16); ok {
		rrType = qType
	}

	for _, rr := range repo.rrs {
		if rr.Type == rrType && rr.Domain == domain && rr.Name == name {
			rrs = append(rrs, rr.RR)
		}
	}

	if len(rrs) == 0 {
		for _, rr := range repo.rrs {
			if rr.Type == rrType && rr.Domain == domain && rr.Name == wildCardName {
				rrs = append(rrs, rr.RR)
			}
		}
	}

	return rrs
}

type resolveName string

func (name resolveName) GetDomain() string {
	v := strings.TrimSuffix(string(name), ".")
	x := strings.Split(v, ".")
	if len(x) < 3 {
		return v
	}

	return strings.Join(x[len(x)-2:], ".")
}

func (name resolveName) GetName() string {
	v := strings.TrimSuffix(string(name), ".")
	x := strings.Split(v, ".")
	if len(x) < 3 {
		return "@"
	}

	return strings.Join(x[:len(x)-2], ".")
}

func (name resolveName) GetWildCardName() string {
	v := strings.TrimSuffix(string(name), ".")
	x := strings.Split(v, ".")
	if len(x) == 3 {
		return "*"
	}

	switch len(x) {
	case 0, 1, 2:
		return ""
	case 3:
		return "*"
	default:
		return "*." + strings.Join(x[1:len(x)-2], ".")
	}
}
