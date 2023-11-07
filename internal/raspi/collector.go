//go:build linux && arm64
// +build linux,arm64

package raspi

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace = "raspi"
)

type Collector struct {
	Up          *prometheus.Desc
	CpuTemp     *prometheus.Desc
	MemoryUsage *prometheus.Desc
	CpuUsage    *prometheus.Desc
	NetReceived *prometheus.Desc
	NetSent     *prometheus.Desc
}

func NewCollector() *Collector {
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname, Exiting application")
	}

	return &Collector{
		Up: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "up"),
			"Current health status of the server (1 = UP, 0 = DOWN).", nil, prometheus.Labels{"host": hostname}),
		CpuTemp: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "cpu_temp_celcius"),
			"CPU temperature", nil, prometheus.Labels{"host": hostname}),
		MemoryUsage: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "memory_usage_mb"),
			"Memory Usage", nil, prometheus.Labels{"host": hostname}),
		CpuUsage: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "cpu_usage_percentage"),
			"CPU Usage", []string{"cpu"}, prometheus.Labels{"host": hostname}),
		NetReceived: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "network_received_kb"),
			"Network Received bytes", []string{"type", "interface"}, prometheus.Labels{"host": hostname}),
		NetSent: prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "network_sent_kb"),
			"Network Sent bytes", []string{"type", "interface"}, prometheus.Labels{"host": hostname}),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.Up
	ch <- collector.CpuTemp
	ch <- collector.MemoryUsage
	ch <- collector.CpuUsage
	ch <- collector.NetReceived
	ch <- collector.NetSent
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Println("Starting Raspi Metrics collection")

	up := prometheus.MustNewConstMetric(collector.Up, prometheus.GaugeValue, 1.0)
	up = prometheus.NewMetricWithTimestamp(time.Now(), up)
	ch <- up

	cpuTemp := getCpuTemp()
	cpuTempMetric := prometheus.MustNewConstMetric(collector.CpuTemp, prometheus.GaugeValue, cpuTemp)
	cpuTempMetric = prometheus.NewMetricWithTimestamp(time.Now(), cpuTempMetric)
	ch <- cpuTempMetric

	memoryUsage := getUsedMemory()
	memoryUsageMetric := prometheus.MustNewConstMetric(collector.MemoryUsage, prometheus.GaugeValue, memoryUsage)
	memoryUsageMetric = prometheus.NewMetricWithTimestamp(time.Now(), memoryUsageMetric)
	ch <- memoryUsageMetric

	cpuUsages := getCpuUsage()
	for i, cpuUsage := range cpuUsages {
		labelValue := fmt.Sprintf("cpu%d", i-1)
		if i == 0 {
			labelValue = "cpu"
		}
		cpuUsageMetric := prometheus.MustNewConstMetric(collector.CpuUsage, prometheus.GaugeValue, cpuUsage, labelValue)
		cpuUsageMetric = prometheus.NewMetricWithTimestamp(time.Now(), cpuUsageMetric)
		ch <- cpuUsageMetric
	}

	transmittedNetworkMetrics := getNetworkMetrics("tx")
	for key, val := range transmittedNetworkMetrics {
		transmittedNetworkBytesMetrics := prometheus.MustNewConstMetric(collector.NetSent, prometheus.GaugeValue,
			float64(val["tx_bytes"]), "tx_bytes", key)
		transmittedNetworkBytesMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedNetworkBytesMetrics)
		ch <- transmittedNetworkBytesMetrics

		transmittedErrorsMetrics := prometheus.MustNewConstMetric(collector.NetSent, prometheus.GaugeValue,
			float64(val["tx_errors"]), "tx_errors", key)
		transmittedErrorsMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedErrorsMetrics)
		ch <- transmittedErrorsMetrics

		transmittedDroppedMetrics := prometheus.MustNewConstMetric(collector.NetSent, prometheus.GaugeValue,
			float64(val["tx_drops"]), "tx_drops", key)
		transmittedDroppedMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedDroppedMetrics)
		ch <- transmittedDroppedMetrics
	}

	receievedNetworkMetrics := getNetworkMetrics("rx")
	for key, val := range receievedNetworkMetrics {
		receivedNetworkBytesMetrics := prometheus.MustNewConstMetric(collector.NetReceived, prometheus.GaugeValue,
			float64(val["rx_bytes"]), "rx_bytes", key)
		receivedNetworkBytesMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedNetworkBytesMetrics)
		ch <- receivedNetworkBytesMetrics

		receivedErrorsMetrics := prometheus.MustNewConstMetric(collector.NetReceived, prometheus.GaugeValue,
			float64(val["rx_errors"]), "rx_errors", key)
		receivedErrorsMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedErrorsMetrics)
		ch <- receivedErrorsMetrics

		receivedDroppedMetrics := prometheus.MustNewConstMetric(collector.NetReceived, prometheus.GaugeValue,
			float64(val["rx_drops"]), "rx_drops", key)
		receivedDroppedMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedDroppedMetrics)
		ch <- receivedDroppedMetrics
	}

	log.Println("Finished Raspi Metrics collection")

}
