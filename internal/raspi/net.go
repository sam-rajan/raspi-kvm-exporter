package raspi

import (
	"fmt"
	"log"
	"os"
	"raspikvm_exporter/internal/config"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func getSentMetrics(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {
	transmittedNetworkMetrics := getNetworkMetrics("tx")
	for key, val := range transmittedNetworkMetrics {
		transmittedNetworkBytesMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["tx_bytes"]), "tx_bytes", key)
		transmittedNetworkBytesMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedNetworkBytesMetrics)
		ch <- transmittedNetworkBytesMetrics

		transmittedErrorsMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["tx_errors"]), "tx_errors", key)
		transmittedErrorsMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedErrorsMetrics)
		ch <- transmittedErrorsMetrics

		transmittedDroppedMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["tx_drops"]), "tx_drops", key)
		transmittedDroppedMetrics = prometheus.NewMetricWithTimestamp(time.Now(), transmittedDroppedMetrics)
		ch <- transmittedDroppedMetrics
	}
}

func getReceievedMetrics(ch chan<- prometheus.Metric, desc *prometheus.Desc, config config.RaspiCollectorConfig) {
	receievedNetworkMetrics := getNetworkMetrics("rx")
	for key, val := range receievedNetworkMetrics {
		receivedNetworkBytesMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["rx_bytes"]), "rx_bytes", key)
		receivedNetworkBytesMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedNetworkBytesMetrics)
		ch <- receivedNetworkBytesMetrics

		receivedErrorsMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["rx_errors"]), "rx_errors", key)
		receivedErrorsMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedErrorsMetrics)
		ch <- receivedErrorsMetrics

		receivedDroppedMetrics := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
			float64(val["rx_drops"]), "rx_drops", key)
		receivedDroppedMetrics = prometheus.NewMetricWithTimestamp(time.Now(), receivedDroppedMetrics)
		ch <- receivedDroppedMetrics
	}

}

func getNetworkMetrics(kind string) map[string]map[string]int64 {

	result := make(map[string]map[string]int64)
	netInterfaces, err := os.ReadDir("/sys/class/net/")

	if err != nil {
		log.Printf("Error getting network transaction. Error: %s", err.Error())
		return result
	}

	for i := 0; i < 2; i++ {

		for _, netInterface := range netInterfaces {

			interfaceMetrics, ok := result[netInterface.Name()]
			if !ok {
				interfaceMetrics = map[string]int64{}
			}

			fileName := fmt.Sprintf("/sys/class/net/%s/statistics/%s_bytes", netInterface.Name(), kind)
			netInfo, err := os.ReadFile(fileName)
			if err != nil {
				log.Printf("Failed to open %s, Error: %s", fileName, err.Error())
				continue
			}

			trafficBytes, _ := strconv.ParseInt(string(netInfo[:len(netInfo)-1]), 10, 64)

			fileName = fmt.Sprintf("/sys/class/net/%s/statistics/%s_errors", netInterface.Name(), kind)
			netErrors, err := os.ReadFile(fileName)
			if err != nil {
				log.Printf("Failed to open %s, Error: %s", fileName, err.Error())
				continue
			}

			trafficErrors, _ := strconv.ParseInt(string(netErrors[:len(netErrors)-1]), 10, 64)

			fileName = fmt.Sprintf("/sys/class/net/%s/statistics/%s_dropped", netInterface.Name(), kind)
			netDropped, err := os.ReadFile(fileName)
			if err != nil {
				log.Printf("Failed to open %s, Error: %s", fileName, err.Error())
				continue
			}

			trafficDropped, _ := strconv.ParseInt(string(netDropped[:len(netDropped)-1]), 10, 64)
			interfaceMetrics[kind+"_bytes"] = abs(toKb(trafficBytes) - interfaceMetrics[kind+"_bytes"])
			interfaceMetrics[kind+"_drops"] = abs(interfaceMetrics[kind+"_drops"] - toKb(trafficDropped))
			interfaceMetrics[kind+"_errors"] = abs(interfaceMetrics[kind+"_errors"] - toKb(trafficErrors))

			result[netInterface.Name()] = interfaceMetrics
		}

		time.Sleep(500 * time.Millisecond)
	}

	return result
}

func abs(val int64) int64 {
	if val < 0 {
		val = val * -1
	}
	return val
}

func toKb(val int64) int64 {
	return (val / 1000) + (val % 1000)
}
