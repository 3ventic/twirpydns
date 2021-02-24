package client

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/3ventic/twirpydns/workers"
	"github.com/cenkalti/backoff/v4"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type Server struct {
	Client          twirpydns.TwirpyDNS
	Secret          string
	FallbackClient  *dns.Client
	FallbackEnabled bool
	FallbackAddress string
	FallbackTimeout time.Duration
	Timeout         time.Duration
	Worker          workers.Worker
}

type results struct {
	twirp    *twirpydns.DNSResponse
	fallback *dns.Msg
	err      error
}

const (
	backoffDuration = 100 * time.Millisecond
)

func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	defer w.Close()

	m, err := r.Pack()
	if err != nil {
		log.Printf("packing: %v", err)
		return
	}

	qs := make([]string, len(r.Question))
	for i := range r.Question {
		qs[i] = r.Question[i].String()
	}

	id := strings.Join(qs, "!!!")

	workerChannel, first := s.Worker.Run(id, s.buildRequester(m, r))
	log.Printf("%s: %v", id, first)

	workerResult := <-workerChannel

	var msg *dns.Msg

	res, ok := workerResult.(*results)

	switch {
	case !ok:
		log.Printf("invalid results of type %T", workerResult)
	case res.err != nil:
		log.Printf("error from worker: %v", err)
	case res.fallback != nil:
		msg = res.fallback
	default:
		msg = new(dns.Msg)

		err := msg.Unpack(res.twirp.Msg)
		if err != nil {
			log.Printf("unpacking twirp response: %v", err)
			return
		}
	}

	// ensure response ID matches request ID
	msg.Id = r.Id
	err = w.WriteMsg(msg)

	if err != nil {
		log.Printf("writing: %v", err)
	}
}

func (s *Server) buildRequester(m []byte, r *dns.Msg) func() interface{} {
	return func() interface{} {
		ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
		defer cancel()

		retryer := backoff.WithContext(backoff.NewConstantBackOff(backoffDuration), ctx)

		var (
			res *twirpydns.DNSResponse
			in  *dns.Msg
		)

		err := backoff.Retry(func() error {
			var err error
			res, err = s.Client.DNS(ctx, &twirpydns.DNSRequest{
				Msg:    m,
				Secret: s.Secret,
			})
			return err
		}, retryer)

		if err != nil {
			if !s.FallbackEnabled {
				err = errors.Wrap(err, "requesting")
			} else {
				log.Printf("requesting: %v", err)

				fbCtx, cancelFallback := context.WithTimeout(context.Background(), s.FallbackTimeout)
				defer cancelFallback()

				// fallback
				in, _, err = s.FallbackClient.ExchangeContext(fbCtx, r, s.FallbackAddress)
				if err != nil {
					err = errors.Wrap(err, "requesting fallback")
				}
			}
		}

		return &results{
			twirp:    res,
			fallback: in,
			err:      err,
		}
	}
}
