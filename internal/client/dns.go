package client

import (
	"context"
	"log"
	"time"

	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/cenkalti/backoff/v4"
	"github.com/miekg/dns"
)

type Server struct {
	Client          twirpydns.TwirpyDNS
	Secret          string
	FallbackAddress string
	Timeout         time.Duration
}

var fallbackClient = &dns.Client{}

func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()
	m, err := r.Pack()
	if err != nil {
		log.Printf("packing: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	retryer := backoff.WithContext(backoff.NewConstantBackOff(100*time.Millisecond), ctx)
	var res *twirpydns.DNSResponse
	err = backoff.Retry(func() error {
		res, err = s.Client.DNS(ctx, &twirpydns.DNSRequest{
			Msg:    m,
			Secret: s.Secret,
		})
		return err
	}, retryer)
	if err != nil {
		log.Printf("requesting: %v", err)

		// fallback
		var in *dns.Msg
		in, _, err = fallbackClient.Exchange(r, s.FallbackAddress)
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
