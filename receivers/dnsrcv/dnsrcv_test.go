package dnsrcv_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ciphermarco/BOAST/log"
	"github.com/ciphermarco/BOAST/receivers/dnsrcv"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

var exampleDomain = "example.com."
var exampleIP = "203.0.113.77"

func NewTestHandler() *dnsrcv.ExportDNSHandler {
	return dnsrcv.NewExportDNSHandler(exampleDomain, exampleIP, []string{"testing"}, &mockStorage{})
}

func TestDNSResponseA(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeA)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		a, ok := qr.Msg.Answer[0].(*dns.A)
		if !ok {
			t.Fatal("wrong type")
		}

		if a.A.String() != exampleIP {
			t.Errorf("wrong A: %v (want) != %v (got)", exampleIP, a.A)
		}
	}
}

func TestDNSResponseNS(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeNS)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		ns, ok := qr.Msg.Answer[0].(*dns.NS)
		if !ok {
			t.Fatal("wrong type")
		}

		want := "ns1." + exampleDomain
		if ns.Ns != want {
			t.Errorf("wrong Ns: %v (want) != %v (got)", want, ns.Ns)
		}

		ns2, ok := qr.Msg.Answer[1].(*dns.NS)
		if !ok {
			t.Fatal("wrong type")
		}

		want2 := "ns2." + exampleDomain
		if ns.Ns != want {
			t.Errorf("wrong Ns: %v (want) != %v (got)", want2, ns2.Ns)
		}
	}
}

func TestDNSResponseSOA(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeSOA)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		soa, ok := qr.Msg.Answer[0].(*dns.SOA)
		if !ok {
			t.Fatal("wrong type")
		}

		wantNs := "ns1." + exampleDomain
		if soa.Ns != wantNs {
			t.Errorf("wrong Ns: %v (want) != %v (got)", wantNs, soa.Ns)
		}

		wantMbox := "mail." + exampleDomain
		if soa.Mbox != wantMbox {
			t.Errorf("wrong Mbox: %v (want) != %v (got)", wantMbox, soa.Mbox)
		}
	}
}

func TestDNSResponseMX(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeMX)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		mx, ok := qr.Msg.Answer[0].(*dns.MX)
		if !ok {
			t.Fatal("wrong type")
		}

		wantMx := "mail." + exampleDomain
		if mx.Mx != wantMx {
			t.Errorf("wrong Mx: %v (want) != %v (got)", wantMx, mx.Mx)
		}
	}
}

func TestDNSResponseTXT(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeTXT)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		txt, ok := qr.Msg.Answer[0].(*dns.TXT)
		if !ok {
			t.Fatal("wrong type")
		}

		wantTxt := []string{"testing"}
		if !reflect.DeepEqual(txt.Txt, wantTxt) {
			t.Errorf("wrong Txt: %v (want) != %v (got)", wantTxt, txt.Txt)
		}
	}
}

func TestDNSResponseANY(t *testing.T) {
	handler := NewTestHandler()

	for _, n := range []string{"example.com.", "sub.example.com."} {
		qr := dnstest.NewRecorder(&test.ResponseWriter{})
		dnsMsg := &dns.Msg{}
		dnsMsg.SetQuestion(n, dns.TypeANY)

		handler.ServeDNS(qr, dnsMsg)

		if qr.Msg == nil {
			t.Fatal("got nil message")
		}

		if qr.Msg.Rcode == dns.RcodeNameError {
			t.Errorf("expected NOERROR got %s", dns.RcodeToString[qr.Msg.Rcode])
		}

		for _, rr := range qr.Msg.Answer {
			switch rr.(type) {
			case *dns.A:
				a := rr.(*dns.A)
				if a.A.String() != exampleIP {
					t.Errorf("wrong A: %v (want) != %v (got)", exampleIP, a.A)
				}
			case *dns.NS:
				ns := rr.(*dns.NS)
				want := "ns1." + exampleDomain
				want2 := "ns2." + exampleDomain
				if ns.Ns != want && ns.Ns != want2 {
					t.Errorf("wrong Ns: %v (want) != %v (got)", want, ns.Ns)
				}
			case *dns.SOA:
				soa := rr.(*dns.SOA)
				wantNs := "ns1." + exampleDomain
				if soa.Ns != wantNs {
					t.Errorf("wrong Ns: %v (want) != %v (got)", wantNs, soa.Ns)
				}

				wantMbox := "mail." + exampleDomain
				if soa.Mbox != wantMbox {
					t.Errorf("wrong Mbox: %v (want) != %v (got)", wantMbox, soa.Mbox)
				}
			case *dns.MX:
				mx := rr.(*dns.MX)
				wantMx := "mail." + exampleDomain
				if mx.Mx != wantMx {
					t.Errorf("wrong Mx: %v (want) != %v (got)", wantMx, mx.Mx)
				}
			case *dns.TXT:
				txt := rr.(*dns.TXT)
				wantTxt := []string{"testing"}
				if !reflect.DeepEqual(txt.Txt, wantTxt) {
					t.Errorf("wrong Txt: %v (want) != %v (got)", wantTxt, txt.Txt)
				}
			}
		}
	}
}
