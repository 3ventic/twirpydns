package main

import (
	"log"
	"net/http"

	"github.com/3ventic/twirpydns/internal/server"
	"github.com/3ventic/twirpydns/rpc/twirpydns"
	"github.com/go-ini/ini"
	"github.com/miekg/dns"
)

func main() {
	cfg, err := ini.Load("server.ini")
	if err != nil {
		log.Fatalf("failed to read server.ini: %v", err)
	}

	server := server.Server{
		Address: cfg.Section("").Key("upstream").MustString("1.1.1.1:53"),
		Secret:  cfg.Section("").Key("secret").String(),
		Client:  &dns.Client{},
	}
	twirpHandler := twirpydns.NewTwirpyDNSServer(server)

	err = http.ListenAndServe(cfg.Section("").Key("listen_address").MustString(":8080"), twirpHandler)
	if err != nil {
		log.Fatal(err)
	}
}
