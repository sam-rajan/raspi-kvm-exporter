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

	virshMetrics = &VirshMetrics{
		VMStatus:     map[string]float64{},
		MemoryStatus: map[string]float64{},
		CPUTime:      map[string]float64{},
	}

	l := libvirt.NewWithDialer(dialers.NewLocal())
	if err := l.Connect(); err != nil {
		log.Printf("Failed to connect: %s", err.Error())
		return
	}

	defer l.Disconnect()
	domains, _, err := l.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive)
	if err != nil {
		log.Printf("failed to retrieve domains: %v", err)
		return
	}

	for _, domain := range domains {
		state, _, mem, _, cpuTime, _ := l.DomainGetInfo(domain)
		virshMetrics.VMStatus[domain.Name] = float64(state)
		memKb := (mem / 1000) + (mem % 1000)
		memMb := (memKb / 1024) + (memKb % 1024)
		virshMetrics.MemoryStatus[domain.Name] = float64(memMb)
		virshMetrics.CPUTime[domain.Name] = float64(cpuTime)
	}

	for _, val := range virshMetrics.VMStatus {
		domainUpMetric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, val)
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
