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

const (
	defaultTimeout = 2 * time.Second
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
		Timeout:         cfg.Section("").Key("timeout").MustDuration(defaultTimeout),
		FallbackEnabled: cfg.Section("fallback").Key("enabled").MustBool(true),
		FallbackAddress: cfg.Section("fallback").Key("address").MustString("1.1.1.1:53"),
		FallbackTimeout: cfg.Section("fallback").Key("timeout").MustDuration(defaultTimeout),
		Worker:          workers.New(),
	}

	err = dns.ListenAndServe("127.0.0.1:53", "udp", s)
	if err != nil {
		log.Fatal(err)
	}
}
