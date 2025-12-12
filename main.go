package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Starting rgw-usage-exporter")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Delay start collectors
	if config.StartDelay > 0 {
		log.Printf("Start delay %d seconds...", config.StartDelay)
		time.Sleep(time.Duration(config.StartDelay) * time.Second)
	}

	// Run collectors metric RGW
	startRGWStatCollector(config)

	// Register exporter in Prometheus
	exporter := NewRGWExporter(config)
	prometheus.MustRegister(exporter)

	// HTTP-handler for /metrics
	http.Handle("/metrics", promhttp.Handler())

	listenAddr := fmt.Sprintf("%s:%d", config.ListenIP, config.ListenPort)
	log.Printf("Serving metrics on http://%s/metrics", listenAddr)

	// Run HTTP-server
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}
