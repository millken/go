package main

import (
	"github.com/miekg/dns"
	"github.com/qiniu/log"
)


func listenAndServe(ip string) {

	prots := []string{"udp", "tcp"}

	for _, prot := range prots {
		go func(p string) {
			server := &dns.Server{Addr: ip, Net: p}

			log.Printf("Opening on %s %s", ip, p)
			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("geodns: failed to setup %s %s: %s", ip, p, err)
			}
			log.Fatalf("mkdns: ListenAndServe unexpectedly returned")
		}(prot)
	}

}
