package raspi

import (
	"bufio"
	"log"
	"os"
	"raspikvm_exporter/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func getUsedMemory(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {

	memInfo, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0444)
	if err != nil {
		log.Println("Failed to open /proc/meminfo, Error: ", err.Error())
		return
	}

	defer memInfo.Close()

	scanner := bufio.NewScanner(memInfo)
	scanner.Split(bufio.ScanLines)

	var totalMemory, freeMemory, buffer, cached int64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			line = strings.ReplaceAll(line, "MemTotal:", "")
			line = strings.ReplaceAll(line, " ", "")
			line = strings.ReplaceAll(line, "kB", "")
			totalMemory, _ = strconv.ParseInt(line, 10, 64)
			continue
		}

		if strings.HasPrefix(line, "MemFree:") {
			line = strings.ReplaceAll(line, "MemFree:", "")
			line = strings.ReplaceAll(line, " ", "")
			line = strings.ReplaceAll(line, "kB", "")
			freeMemory, _ = strconv.ParseInt(line, 10, 64)
			continue
		}

		if strings.HasPrefix(line, "Buffers:") {
			line = strings.ReplaceAll(line, "Buffers:", "")
			line = strings.ReplaceAll(line, " ", "")
			line = strings.ReplaceAll(line, "kB", "")
			buffer, _ = strconv.ParseInt(line, 10, 64)
			continue
		}

		if strings.HasPrefix(line, "Cached:") {
			line = strings.ReplaceAll(line, "Cached:", "")
			line = strings.ReplaceAll(line, " ", "")
			line = strings.ReplaceAll(line, "kB", "")
			cached, _ = strconv.ParseInt(line, 10, 64)
			continue
		}
	}

	usedMemory := totalMemory - freeMemory - (buffer + cached)
	usedMemory = (usedMemory / 1024) + (usedMemory % 1024)
	memoryUsageMetric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(usedMemory))
	memoryUsageMetric = prometheus.NewMetricWithTimestamp(time.Now(), memoryUsageMetric)
	ch <- memoryUsageMetric
}
