//go:build (linux && ignore) || arm64
// +build linux,ignore arm64

package raspi

import (
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace = "raspi_kvm"
)

type Collector struct {
	Up      *prometheus.Desc
	CpuTemp *prometheus.Desc
}

func NewCollector() *Collector {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname, Exiting application")
	}

	return &Collector{
		Up: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "up"),
			"Current health status of the server (1 = UP, 0 = DOWN).", nil, prometheus.Labels{"host": hostname}),
		CpuTemp: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "raspi", "cpu_temp_celcius"),
			"CPU temperature", nil, prometheus.Labels{"host": hostname}),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.Up
	ch <- collector.CpuTemp
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	err, cpuTemp := getCpuTemp()

	if err != nil {
		return
	}
	cpuTempMetric := prometheus.MustNewConstMetric(collector.CpuTemp, prometheus.GaugeValue, cpuTemp)
	cpuTempMetric = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), cpuTempMetric)

	up := prometheus.MustNewConstMetric(collector.Up, prometheus.GaugeValue, 1.0)
	up = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), up)

	ch <- up
	ch <- cpuTempMetric
}
