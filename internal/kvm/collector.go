//go:build linux && arm64
// +build linux,arm64

package kvm

import (
	"log"
	"os"
	"raspikvm_exporter/internal/config"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace = "kvm"
)

var vmCollectorConfig *config.KvmCollectorConfig

var metricGenerator = map[string]func(chan<- prometheus.Metric, *prometheus.Desc, config.KvmCollectorConfig){
	"domainsUp":         getDomainsUp,
	"domainMemoryUsage": getDomainMemoryUsage,
	"domainCpuUsage":    getDomainCpuUsage,
}

var syncTask = map[string]bool{
	"domainsUp": true,
}

type Collector struct {
	Metrices map[string]*prometheus.Desc
}

func NewCollector(collectorConfig *config.KvmCollectorConfig) *Collector {
	vmCollectorConfig = collectorConfig
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname, Exiting application")
	}

	return &Collector{
		Metrices: map[string]*prometheus.Desc{
			"domainsUp": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domains_up"),
				"Number of domains up", []string{"vm"}, prometheus.Labels{"host": hostname}),
			"domainMemoryUsage": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domain_memory_usage_mb"),
				"Virtual machine memory usage", []string{"vm"}, prometheus.Labels{"host": hostname}),
			"domainCpuUsage": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domain_cpu_time_ms"),
				"Virtual machine cpu time", []string{"vm"}, prometheus.Labels{"host": hostname}),
		},
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, val := range collector.Metrices {
		ch <- val
	}
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	log.Println("Starting KVM metrics collection")
	for key, desc := range collector.Metrices {

		log.Println("Collecting metric: ", key)
		if syncTask[key] {
			metricGenerator[key](ch, desc, *vmCollectorConfig)
			continue
		}
		wg.Add(1)
		go func(k string, d *prometheus.Desc) {
			metricGenerator[k](ch, d, *vmCollectorConfig)
			wg.Done()
		}(key, desc)
	}

	wg.Wait()

	log.Println("Finished KVM Metrics collection")
}
