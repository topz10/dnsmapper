package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
	"strings"
)

var VERSION string = "0.0.1"

var (
	flagdomain   = flag.String("domain", "example.com", "base domain for the dnsmapper")
	flagip       = flag.String("ip", "127.0.0.1", "set the IP address")
	flagdnsport  = flag.String("dnsport", "53", "Set the DNS port")
	flaghttpport = flag.String("httpport", "80", "Set the HTTP port")
	flaglog      = flag.Bool("log", false, "be more verbose")
)

var baseLength int

func getUuidFromDomain(name string) string {
	lx := dns.SplitLabels(name)
	ql := lx[0 : len(lx)-baseLength]
	return strings.ToLower(strings.Join(ql, "."))
}

func main() {
	flag.Parse()
	log.Printf("Starting dnsmapper %s\n", VERSION)

	baseLength = dns.LenLabels(*flagdomain)

	log.Println("Listening for requests to", *flagdomain)

	dns.HandleFunc(*flagdomain, setupServerFunc())

	redisConnect()

	go httpHandler()
	go listenAndServeDNS(*flagip + ":" + *flagdnsport)

	terminate := make(chan os.Signal)
	signal.Notify(terminate, os.Interrupt)

	<-terminate
	log.Printf("dnsmapper: signal received, stopping")

}