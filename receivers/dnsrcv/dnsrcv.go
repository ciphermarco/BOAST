package dnsrcv

import (
	"fmt"
	"net"
	"strings"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"

	"github.com/miekg/dns"
)

const shortTTL = 300

// Receiver represents the DNS protocol receiver.
type Receiver struct {
	Name     string
	Domain   string
	Host     string
	Ports    []int
	PublicIP string
	Txt      []string
	Storage  app.Storage
}

// ListenAndServe sets the necessary conditions for the underlying dns.Server
// to serve the BOAST's custom DNS server for each configured port.
//
// For full functionality, this server must be used as nameserver for the domain.
//
// Any errors are returned via the received channel.
func (r *Receiver) ListenAndServe(err chan error) {
	for _, port := range r.Ports {
		go func(p int) {
			addr := r.Host + fmt.Sprintf(":%d", p)
			srv := &dns.Server{
				Addr: addr,
				Net:  "udp",
			}

			srv.Handler = &dnsHandler{
				domain:   r.Domain,
				publicIP: r.PublicIP,
				txt:      r.Txt,
				storage:  r.Storage,
			}

			log.Info("%s: Listening on %s\n", r.Name, addr)
			err <- srv.ListenAndServe()
		}(port)
	}
}

type dnsHandler struct {
	domain   string
	publicIP string
	txt      []string
	storage  app.Storage
}

var queryTypeNames = map[uint16]string{
	dns.TypeA:     "A",
	dns.TypeNS:    "NS",
	dns.TypeSOA:   "SOA",
	dns.TypeMX:    "MX",
	dns.TypeCNAME: "CNAME",
	dns.TypeAAAA:  "AAAA",
	dns.TypeTXT:   "TXT",
}

// ServeDNS is the handler for BOAST's DNS queries.
// It responds to A, NS, SOA, and MX queries always pointing to the same IP.
func (d *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Info("DNS event received")
	msg := dns.Msg{}
	msg.SetReply(r)

	id, canary := d.storage.SearchTest(
		func(key, value string) bool {
			return strings.Contains(msg.Question[0].Name, key)
		},
	)

	if id != "" {
		qTypeName := queryTypeNames[r.Question[0].Qtype]
		evt, err := app.NewDNSEvent(
			id,
			"DNS",
			w.RemoteAddr().String(),
			r.String(),
			qTypeName,
		)
		if err != nil {
			log.Info("Error creating a new DNS event")
			log.Debug("New DNS event error: %v", err)
		} else {
			if err := d.storage.StoreEvent(evt); err != nil {
				log.Info("Error storing a new DNS event")
				log.Debug("Store DNS event error: %v", err)
			} else {
				log.Info("New DNS event stored")
			}
			log.Debug("DNS event object:\n%s", evt.String())
		}
	} else {
		log.Debug("DNS event test not found: id=\"%s\" canary=\"%s\"",
			id, canary)
	}

	d.setDNSAnswer(&msg, r)
	w.WriteMsg(&msg)
}

func (d *dnsHandler) setDNSAnswer(msg, r *dns.Msg) {
	qName := msg.Question[0].Name
	if strings.HasSuffix(toFQDN(qName), toFQDN(d.domain)) {
		msg.Authoritative = true
		hdr := dns.RR_Header{
			Name:  qName,
			Class: dns.ClassINET,
			Ttl:   shortTTL,
		}

		qType := r.Question[0].Qtype

		if qType == dns.TypeA || qType == dns.TypeANY {
			hdr.Rrtype = dns.TypeA
			msg.Answer = append(msg.Answer,
				&dns.A{
					Hdr: hdr,
					A:   net.ParseIP(d.publicIP),
				})
		}

		if qType == dns.TypeNS || qType == dns.TypeANY {
			hdr.Rrtype = dns.TypeNS
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: hdr,
				Ns:  "ns1." + toFQDN(d.domain),
			})
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: hdr,
				Ns:  "ns2." + toFQDN(d.domain),
			})
		}

		if qType == dns.TypeSOA || qType == dns.TypeANY {
			hdr.Rrtype = dns.TypeSOA
			msg.Answer = append(msg.Answer, &dns.SOA{
				Hdr:     hdr,
				Ns:      "ns1." + toFQDN(d.domain),
				Mbox:    "mail." + toFQDN(d.domain),
				Refresh: 604800,
				Serial:  10000,
				Retry:   11000,
				Expire:  120000,
				Minttl:  10000,
			})
		}

		if qType == dns.TypeMX || qType == dns.TypeANY {
			hdr.Rrtype = dns.TypeMX
			msg.Answer = append(msg.Answer, &dns.MX{
				Hdr:        hdr,
				Preference: 1,
				Mx:         "mail." + toFQDN(d.domain),
			})
		}

		if len(d.txt) > 0 {
			if qType == dns.TypeTXT || qType == dns.TypeANY {
				hdr.Rrtype = dns.TypeTXT
				msg.Answer = append(msg.Answer, &dns.TXT{
					Hdr: hdr,
					Txt: d.txt,
				})
			}
		}
	}
}

func toFQDN(s string) string {
	l := len(s)
	if l == 0 || s[l-1] == '.' {
		return strings.ToLower(s)
	}
	return strings.ToLower(s) + "."
}
