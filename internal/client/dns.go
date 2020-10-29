package client

import (
	"context"
	"log"

	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/miekg/dns"
)

type Server struct {
	Client twirpydns.TwirpyDNS
	Secret string
}

var fallbackClient = &dns.Client{}

func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	m, err := r.Pack()
	if err != nil {
		log.Printf("packing: %v", err)
		return
	}

	res, err := s.Client.DNS(context.Background(), &twirpydns.DNSRequest{
		Msg:    m,
		Secret: s.Secret,
	})
	if err != nil {
		log.Printf("requesting: %v", err)

		// fallback
		var in *dns.Msg
		in, _, err = fallbackClient.Exchange(r, "1.1.1.1:53")
		if err != nil {
			log.Printf("requesting fallback: %v", err)
			return
		}

		err = w.WriteMsg(in)
		if err != nil {
			log.Printf("writing fallback: %v", err)
			return
		}
		return
	}

	_, err = w.Write(res.Msg)
	if err != nil {
		log.Printf("writing: %v", err)
		return
	}
}
