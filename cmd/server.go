package main

import (
	"log"
	"net/http"
	"raspikvm_exporter/internal/kvm"
	"raspikvm_exporter/internal/raspi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	raspiCollector := raspi.NewCollector()
	kvmCollector := kvm.NewCollector()
	prometheus.MustRegister(raspiCollector, kvmCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter starting at port 8082")
	err := http.ListenAndServe(":8082", nil)

	if err != nil {
		log.Fatal(err.Error())
	}

}
