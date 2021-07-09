package fns

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type FNS struct {
	Next    plugin.Handler
	RRsRepo RRsRepository
}

func (f *FNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	m.Authoritative = true

	m.Answer = f.Resolve(state.QName(), state.QType(), state.QClass())

	if len(m.Answer) == 0 && f.Next != nil {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}

	state.SizeAndDo(m)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func (f FNS) Resolve(name string, rrType, rrClass uint16) []dns.RR {
	var (
		rrs     []RR
		answers = make([]dns.RR, 0)
	)

	rrs = f.RRsRepo.GetRRs(Condition{
		"name": name,
		"type": dns.TypeCNAME,
	})

	for _, rr := range rrs {
		r := new(dns.CNAME)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: rr.TTL}
		r.Target = strings.TrimRight(rr.Value, ".") + "."
		answers = append(answers, r)
		return answers
	}

	rrs = f.RRsRepo.GetRRs(Condition{
		"name": name,
		"type": rrType,
	})

	for _, rr := range rrs {
		switch rrType {
		case dns.TypeA:
			r := new(dns.A)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: rr.TTL}
			if r.A = net.ParseIP(rr.Value).To4(); r.A != nil {
				answers = append(answers, r)
			}
		case dns.TypeNS:
			r := new(dns.NS)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: rr.TTL}
			r.Ns = strings.TrimRight(rr.Value, ".") + "."
			answers = append(answers, r)
		case dns.TypeAAAA:
			r := new(dns.AAAA)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: rr.TTL}
			if r.AAAA = net.ParseIP(rr.Value).To16(); r.AAAA != nil {
				answers = append(answers, r)
			}
		case dns.TypeMX:
			r := new(dns.MX)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: rr.TTL}
			mx := rr.ValueToMX()
			r.Mx = strings.TrimRight(mx.Domain, ".") + "."
			r.Preference = mx.Priority
			answers = append(answers, r)
		case dns.TypeTXT:
			r := new(dns.TXT)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: rr.TTL}
			r.Txt = strings.Split(rr.Value, "\n")
			answers = append(answers, r)
		case dns.TypeSRV:
			r := new(dns.SRV)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: rr.TTL}
			srv := rr.ValueToSRV()
			if srv.Target != "" {
				r.Priority = srv.Priority
				r.Weight = srv.Weight
				r.Port = srv.Port
				r.Target = strings.TrimRight(srv.Target, ".") + "."
				answers = append(answers, r)
			}
		case dns.TypeCAA:
			r := new(dns.CAA)
			r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeCAA, Class: dns.ClassINET, Ttl: rr.TTL}
			caa := rr.ValueToCAA()
			if caa.Value != "" {
				r.Tag = caa.Tag
				r.Flag = caa.Flags
				r.Value = caa.Value
				answers = append(answers, r)
			}
		}
	}

	return answers
}

func (f FNS) Name() string { return PluginName }
