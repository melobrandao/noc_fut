package main

import (
	"automacao-api-go/internal/dashboard"
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "Endereço local do dashboard")
	flag.Parse()
	handler, err := dashboard.New()
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Addr: *addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("NOC FUT • http://%s • dados demonstrativos", *addr)
	log.Fatal(server.ListenAndServe())
}
