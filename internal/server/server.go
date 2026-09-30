package server

import (
	"context"
	"log"
	"net/netip"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
	"github.com/marttimourujarvi/gateway-dns/internal/dnsstore"
)

type handler struct {
	store *dnsstore.DNSStore
}

func (h *handler) ServeDNS(c context.Context, w dns.ResponseWriter, r *dns.Msg) {
	defer log.Println("ServeDNS done")
	reply := r.Copy()

	reply.Response = true

	for _, q := range r.Question {
		if a, ok := q.(*dns.A); ok {
			// it's an A record question
			log.Println("its A record", q.Header().Name)

			ip, ok := h.store.Lookup(q.Header().Name)

			if !ok {
				reply.Rcode = dns.RcodeNameError
				reply.Pack()
				reply.WriteTo(w)
				return
			}

			rr := &dns.A{
				Hdr: dns.Header{
					Name:  a.Hdr.Name,
					Class: dns.ClassINET,
					TTL:   60,
				},
				A: rdata.A{Addr: netip.MustParseAddr(ip)},
			}
			reply.Answer = append(reply.Answer, rr)
		}

	}

	reply.Pack()
	reply.WriteTo(w)
}

// New builds a UDP DNS server that answers A queries from the shared store.
func New(addr string, store *dnsstore.DNSStore) *dns.Server {
	return &dns.Server{
		Addr:    addr,
		Net:     "udp",
		Handler: &handler{store: store},
	}
}
