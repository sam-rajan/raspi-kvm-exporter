package main

import (
	"flag"
	"log"
	"net/http"
	"raspikvm_exporter/internal/config"
	"raspikvm_exporter/internal/kvm"
	"raspikvm_exporter/internal/raspi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenPort := flag.String("port", "9000", "Port to used by server to listen.")
	configFile := flag.String("config", "", "Path to config file.")
	flag.Parse()

	collectorConfig := config.NewCollectorConfig(configFile)
	if listenPort != nil {
		collectorConfig.Port = *listenPort
	}

	if *collectorConfig.Collectors.Kvm.Enabled {
		kvmCollector := kvm.NewCollector(&collectorConfig.Collectors.Kvm)
		prometheus.MustRegister(kvmCollector)
	}

	if *collectorConfig.Collectors.Raspi.Enabled {
		raspiCollector := raspi.NewCollector(&collectorConfig.Collectors.Raspi)
		prometheus.MustRegister(raspiCollector)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Exporter starting at port %s", collectorConfig.Port)
	err := http.ListenAndServe(":"+collectorConfig.Port, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

}
