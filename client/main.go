package main

import (
	"log"
	"net/http"

	"github.com/3ventic/twirpydns/internal/client"
	"github.com/3ventic/twirpydns/rpc/twirpydns"
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
		Client: twirpydnsClient,
		Secret: cfg.Section("").Key("secret").String(),
	}

	err = dns.ListenAndServe("127.0.0.1:53", "udp", s)
	if err != nil {
		log.Fatal(err)
	}
}
