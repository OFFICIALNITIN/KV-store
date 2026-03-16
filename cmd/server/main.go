package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	config "github.com/OFFICIALNITIN/KV-store"
	"github.com/OFFICIALNITIN/KV-store/internal/server"
	"github.com/OFFICIALNITIN/KV-store/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, _ := config.LoadConfig("config.yaml")

	kv := store.New(time.Duration(cfg.Storage.CleanupInterval) * time.Second)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Critical: Could not start TCP listener: %v", err)
	}

	defer listener.Close()

	log.Printf("Info: TCP listener started on %s", addr)
	log.Println("Concurrency model: Goroutine per connection")

	go func() {
		metricAddr := fmt.Sprintf(":%d", cfg.Server.MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Info: Metrics available at http://localhost:%d/metrics", cfg.Server.MetricsPort)
		http.ListenAndServe(metricAddr, nil)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}

		go server.HandleConnection(conn, kv)
	}
}
