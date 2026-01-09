package raspi

import (
	"log"
	"os"
	"time"

	"raspikvm_exporter/internal/config"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace = "raspi"
)

var raspiCollectorConfig *config.RaspiCollectorConfig

var metricGenerator = map[string]func(chan<- prometheus.Metric, *prometheus.Desc, config.RaspiCollectorConfig){
	"cpuTemp":     getCpuTemp,
	"cpuUsage":    getCpuUsage,
	"memoryUsage": getUsedMemory,
	"netSent":     getSentMetrics,
	"netReceived": getReceievedMetrics,
	"diskUsage":   getDiskUsage,
	"readBytes":   getDiskReads,
	"writeBytes":  getDiskWrites,
}

type Collector struct {
	Metrices map[string]*prometheus.Desc
}

func NewCollector(config *config.RaspiCollectorConfig) *Collector {
	raspiCollectorConfig = config
	hostname, err := os.Hostname()

	if err != nil {
		log.Fatal("Failed to get hostname, Exiting application")
	}

	return &Collector{

		Metrices: map[string]*prometheus.Desc{
			"up": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "up"),
				"Current health status of the server (1 = UP, 0 = DOWN).", nil, prometheus.Labels{"host": hostname}),
			"cpuTemp": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "cpu_temp_celcius"),
				"CPU temperature", nil, prometheus.Labels{"host": hostname}),
			"memoryUsage": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "memory_usage_mb"),
				"Memory Usage", nil, prometheus.Labels{"host": hostname}),
			"cpuUsage": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "cpu_usage_percentage"),
				"CPU Usage", []string{"cpu"}, prometheus.Labels{"host": hostname}),
			"netReceived": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "network_received_kb"),
				"Network Received bytes", []string{"type", "interface"}, prometheus.Labels{"host": hostname}),
			"netSent": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "network_sent_kb"),
				"Network Sent bytes", []string{"type", "interface"}, prometheus.Labels{"host": hostname}),
			"diskUsage": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "disk_usage_bytes"),
				"Disk Usage", []string{"device", "total", "mountpoint"}, prometheus.Labels{}),
			"readBytes": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "disk_read_bytes"),
				"Disk Read Bytes", []string{"device", "mountpoint"}, prometheus.Labels{}),
			"writeBytes": prometheus.NewDesc(prometheus.BuildFQName(NameSpace, "", "disk_write_bytes"),
				"Disk Write Bytes", []string{"device", "mountpoint"}, prometheus.Labels{}),
		},
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, val := range collector.Metrices {
		ch <- val
	}
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Println("Starting Raspi Metrics collection")

	up := prometheus.MustNewConstMetric(collector.Metrices["up"], prometheus.GaugeValue, 1.0)
	up = prometheus.NewMetricWithTimestamp(time.Now(), up)
	ch <- up

	for key, desc := range collector.Metrices {
		if key == "up" {
			continue
		}

		metricGenerator[key](ch, desc, *raspiCollectorConfig)
	}

	log.Println("Finished Raspi Metrics collection")

}
