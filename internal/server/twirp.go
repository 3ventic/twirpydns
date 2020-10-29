package server

import (
	"context"
	"log"
	"time"

	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/miekg/dns"
	"github.com/twitchtv/twirp"
)

// Server implements the twirpydns twirp-DNS server
type Server struct {
	Secret  string
	Client  *dns.Client
	Address string
}

func (s Server) DNS(ctx context.Context, req *twirpydns.DNSRequest) (*twirpydns.DNSResponse, error) {
	if s.Secret != req.Secret {
		return nil, twirp.InvalidArgumentError("secret", "secret does not match")
	}

	r := new(dns.Msg)
	err := r.Unpack(req.Msg)
	if err != nil {
		return nil, twirp.InvalidArgumentError("msg", "unable to unpack msg. Is it a valid DNS query?")
	}

	if r == nil || len(r.Question) == 0 {
		return nil, twirp.InvalidArgumentError("msg", "empty query in msg")
	}

	in, rtt, err := s.Client.ExchangeContext(ctx, r, s.Address)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	if rtt > time.Second {
		log.Printf("query took %v: %v", rtt, r.Question[0].String())
	}

	out, err := in.Pack()
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &twirpydns.DNSResponse{
		Msg: out,
	}, nil
}
