package main

import (
	"log"
	"net/http"
	"raspikvm_exporter/internal/raspi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	raspiCollector := raspi.NewCollector()
	prometheus.MustRegister(raspiCollector)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8082", nil)

	log.Println("Started RasPI KVM exporter on port 8082")
}
