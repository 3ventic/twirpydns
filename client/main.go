package main

import (
	"log"
	"net/http"
	"time"

	"github.com/3ventic/twirpydns/internal/client"
	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/3ventic/twirpydns/workers"
	"github.com/go-ini/ini"
	"github.com/miekg/dns"
)

func main() {
	cfg, err := ini.Load("client.ini")
	if err != nil {
		log.Fatalf("failed to read client.ini: %v", err)
	}

	twirpydnsClient := twirpydns.NewTwirpyDNSProtobufClient(cfg.Section("").Key("server_address").String(), &http.Client{})

	s := &client.Server{
		Client:          twirpydnsClient,
		Secret:          cfg.Section("").Key("secret").String(),
		FallbackAddress: cfg.Section("").Key("fallback_dns").MustString("1.1.1.1:53"),
		Timeout:         cfg.Section("").Key("timeout").MustDuration(10 * time.Second),
		Worker:          workers.New(),
	}

	err = dns.ListenAndServe("127.0.0.1:53", "udp", s)
	if err != nil {
		log.Fatal(err)
	}
}
