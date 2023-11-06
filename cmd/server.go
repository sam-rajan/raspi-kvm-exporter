package main

import (
	"flag"
	"log"
	"net/http"
	"raspikvm_exporter/internal/kvm"
	"raspikvm_exporter/internal/raspi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenPort := flag.String("exporter.port", ":9000", "Port to used by server to listen.")
	flag.Parse()

	raspiCollector := raspi.NewCollector()
	kvmCollector := kvm.NewCollector()
	prometheus.MustRegister(raspiCollector, kvmCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Exporter starting at port %s", *listenPort)
	err := http.ListenAndServe(":"+*listenPort, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

}
