//go:build linux && arm64
// +build linux,arm64

package kvm

import (
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace = "kvm"
)

type Collector struct {
	DomainsUp         *prometheus.Desc
	DomainMemoryUsage *prometheus.Desc
	DomainCpuUsage    *prometheus.Desc
}

func NewCollector() *Collector {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname, Exiting application")
	}

	return &Collector{
		DomainsUp: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domain_up"),
			"KVM virtual machine status (1 = UP, 0 = DOWN).", []string{"vm"}, prometheus.Labels{"host": hostname}),
		DomainMemoryUsage: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domain_memory_usage_mb"),
			"Virtual machine memory usage", []string{"vm"}, prometheus.Labels{"host": hostname}),
		DomainCpuUsage: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "domain_cpu_time_ms"),
			"Virtual machine cpu time", []string{"vm"}, prometheus.Labels{"host": hostname}),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DomainsUp
	ch <- collector.DomainMemoryUsage
	ch <- collector.DomainCpuUsage
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Println("Starting KVM metrics collection")

	domainMetrics := getVirshMetrics()

	for key, val := range domainMetrics.VMStatus {
		domainUpMetric := prometheus.MustNewConstMetric(collector.DomainsUp, prometheus.GaugeValue, val, key)
		domainUpMetric = prometheus.NewMetricWithTimestamp(time.Now(), domainUpMetric)
		ch <- domainUpMetric
	}

	for key, val := range domainMetrics.MemoryStatus {
		domainMemoryMetric := prometheus.MustNewConstMetric(collector.DomainMemoryUsage, prometheus.GaugeValue, val, key)
		domainMemoryMetric = prometheus.NewMetricWithTimestamp(time.Now(), domainMemoryMetric)
		ch <- domainMemoryMetric
	}

	for key, val := range domainMetrics.CPUTime {
		cpuTimeMetric := prometheus.MustNewConstMetric(collector.DomainCpuUsage, prometheus.GaugeValue, val, key)
		cpuTimeMetric = prometheus.NewMetricWithTimestamp(time.Now(), cpuTimeMetric)
		ch <- cpuTimeMetric
	}

	log.Println("Finished KVM Metrics collection")
}
