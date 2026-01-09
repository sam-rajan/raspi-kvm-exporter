package kvm

import (
	"log"
	"raspikvm_exporter/internal/config"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/prometheus/client_golang/prometheus"
)

type VirshMetrics struct {
	VMStatus     map[string]float64
	MemoryStatus map[string]float64
	CPUTime      map[string]float64
}

var virshMetrics = &VirshMetrics{
	VMStatus:     map[string]float64{},
	MemoryStatus: map[string]float64{},
	CPUTime:      map[string]float64{},
}

func getDomainsUp(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.KvmCollectorConfig) {

	metrics := &VirshMetrics{
		VMStatus:     map[string]float64{},
		MemoryStatus: map[string]float64{},
		CPUTime:      map[string]float64{},
	}

	l := libvirt.NewWithDialer(dialers.NewLocal())
	if err := l.Connect(); err != nil {
		log.Printf("Failed to connect: %s", err.Error())
		return
	}

	domains, _, err := l.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive)
	if err != nil {
		log.Printf("failed to retrieve domains: %v", err)
		return
	}

	for _, domain := range domains {
		state, _, mem, _, cpuTime, _ := l.DomainGetInfo(domain)
		metrics.VMStatus[domain.Name] = float64(state)
		memKb := (mem / 1000) + (mem % 1000)
		memMb := (memKb / 1024) + (memKb % 1024)
		metrics.MemoryStatus[domain.Name] = float64(memMb)
		metrics.CPUTime[domain.Name] = float64(cpuTime)
	}

	for key, val := range virshMetrics.VMStatus {
		domainUpMetric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, key)
		domainUpMetric = prometheus.NewMetricWithTimestamp(time.Now(), domainUpMetric)
		ch <- domainUpMetric
	}

}

func getDomainMemoryUsage(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.KvmCollectorConfig) {
	for key, val := range virshMetrics.MemoryStatus {
		domainMemoryMetric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, key)
		domainMemoryMetric = prometheus.NewMetricWithTimestamp(time.Now(), domainMemoryMetric)
		ch <- domainMemoryMetric
	}
}

func getDomainCpuUsage(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.KvmCollectorConfig) {
	for key, val := range virshMetrics.CPUTime {
		cpuTimeMetric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val, key)
		cpuTimeMetric = prometheus.NewMetricWithTimestamp(time.Now(), cpuTimeMetric)
		ch <- cpuTimeMetric
	}
}
