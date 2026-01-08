package raspi

import (
	"log"
	"path/filepath"

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

func getDiskUsage() []*diskMetrics {
	metrices := []*diskMetrics{}
	parts, _ := disk.Partitions(false)
	for _, part := range parts {
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

	return metrices
}
