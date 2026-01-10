package raspi

import (
	"log"
	"path/filepath"
	"raspikvm_exporter/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/disk"
)

type diskMetrics struct {
	deviceName string
	mountpoint string
	totalSize  uint64
	usedSize   uint64
	readBytes  uint64
	writeBytes uint64
}

var metrices = []*diskMetrics{}

func getDiskUsage(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {
	// clear the metrices
	metrices = []*diskMetrics{}

	if disk, ok := config.Disk["enabled"]; ok && disk == false {
		return
	}

	var deviceList []string
	if devices, ok := config.Disk["devices"]; ok {
		deviceList = devices.([]string)
	}

	parts, _ := disk.Partitions(false)
	for _, part := range parts {

		flag := false
		for _, device := range deviceList {
			if strings.Contains(part.Device, device) {
				flag = true
				break
			}
		}

		if !flag && len(deviceList) > 0 {
			continue
		}

		usage, err := disk.Usage(part.Mountpoint)

		if err != nil {
			log.Println("Failed to get disk usage ", err.Error())
			continue
		}
		metric := &diskMetrics{
			deviceName: part.Device,
			mountpoint: part.Mountpoint,
			totalSize:  usage.Total,
			usedSize:   usage.Used,
		}

		ioCounters, err := disk.IOCounters(metric.deviceName)
		if err == nil {
			stat, ok := ioCounters[filepath.Base(metric.deviceName)]
			if ok {
				metric.readBytes = stat.ReadBytes
				metric.writeBytes = stat.WriteBytes
			}
		}

		metrices = append(metrices, metric)
	}

	for _, diskMetrics := range metrices {
		diskUsage := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(diskMetrics.usedSize), diskMetrics.deviceName,
			strconv.FormatUint(diskMetrics.totalSize, 10), diskMetrics.mountpoint)
		diskUsage = prometheus.NewMetricWithTimestamp(time.Now(), diskUsage)
		ch <- diskUsage
	}
}

func getDiskWrites(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {

	for _, diskMetrics := range metrices {
		writeBytes := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(diskMetrics.writeBytes), diskMetrics.deviceName, diskMetrics.mountpoint)
		writeBytes = prometheus.NewMetricWithTimestamp(time.Now(), writeBytes)
		ch <- writeBytes
	}

}

func getDiskReads(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {
	for _, diskMetrics := range metrices {
		readBytes := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(diskMetrics.readBytes), diskMetrics.deviceName, diskMetrics.mountpoint)
		readBytes = prometheus.NewMetricWithTimestamp(time.Now(), readBytes)
		ch <- readBytes
	}
}
