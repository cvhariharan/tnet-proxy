package main

import (
	"flag"
	"log"
	"os"

	proxy "github.com/cvhariharan/tnet-proxy"
)

func main() {
	listenPort := flag.Int("port", 8000, "port that the proxy listens on")
	target := flag.String("target", "", "target to proxy (eg: localhost:8080, :8080)")
	hostname := flag.String("hostname", "tnet-proxy", "hostname used on the tailnet")
	proto := flag.String("proto", "tcp", "protocol tcp/udp")

	if *proto != "tcp" && *proto != "udp" {
		log.Fatal("protocol should either be tcp or udp")
	}

	flag.Parse()

	if len(os.Getenv("TS_AUTHKEY")) == 0 {
		log.Fatal("environment variable TS_AUTHKEY not set")
	}

	f := proxy.NewForwardProxy(*listenPort, *target, *hostname, *proto)
	log.Fatal(f.ListenAndServe())
}
