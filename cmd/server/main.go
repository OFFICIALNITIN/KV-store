package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/OFFICIALNITIN/KV-store"
	"github.com/OFFICIALNITIN/KV-store/internal/server"
	"github.com/OFFICIALNITIN/KV-store/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, _ := config.LoadConfig("config.yaml")

	kv := store.New(time.Duration(cfg.Storage.CleanupInterval) * time.Second)

	if cfg.Storage.SnapshotFile != "" {
		if err := kv.LoadSnapshot(cfg.Storage.SnapshotFile); err != nil {
			log.Printf("Warning: Failed to load snapshot: %v", err)
		} else {
			log.Println("Info: Successfully loaded data from snapshot!")
		}
	}

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

	go func() {
		httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort)
		mux := http.NewServeMux()
		mux.Handle("/keys/", server.NewHTTPHandler(kv))

		log.Printf("Info: HTTP API listing on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
					break
				}
				log.Printf("Connection error: %v", err)
				continue
			}
			go server.HandleConnection(conn, kv)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("\nInfo: Shutting down server gracefully...")

	if cfg.Storage.SnapshotFile != "" {
		log.Println("Info: Saving memory snapshot to disk...")
		if err := kv.SaveSnapshot(cfg.Storage.SnapshotFile); err != nil {
			log.Printf("Error saving snapshot: %v", err)
		} else {
			log.Println("Info: Snapshot saved successfully!")
		}
	}
}
