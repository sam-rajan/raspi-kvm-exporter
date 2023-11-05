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
	err := http.ListenAndServe(":8082", nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Exporter started at port 8082")
}
