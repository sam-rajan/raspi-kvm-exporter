package kvm

import (
	"log"
	"net"
	"time"

	"github.com/digitalocean/go-libvirt"
)

type VirshMetrics struct {
	VMStatus     map[string]float64
	MemoryStatus map[string]float64
	CPUTime      map[string]float64
}

func getVirshMetrics() VirshMetrics {

	metrics := VirshMetrics{
		VMStatus:     map[string]float64{},
		MemoryStatus: map[string]float64{},
		CPUTime:      map[string]float64{},
	}

	connection, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", time.Duration(2*time.Second))
	if err != nil {
		log.Printf("Failed to dial virsh socket Error: %s", err.Error())
		return metrics
	}
	defer connection.Close()

	l := libvirt.New(connection)
	if err := l.Connect(); err != nil {
		log.Printf("Failed to connect: %s", err.Error())
		return metrics
	}

	domains, err := l.Domains()
	if err != nil {
		log.Printf("failed to retrieve domains: %v", err)
		return metrics
	}

	for _, domain := range domains {
		state, _, mem, _, cpuTime, _ := l.DomainGetInfo(domain)
		metrics.VMStatus[domain.Name] = float64(state)
		metrics.MemoryStatus[domain.Name] = float64((mem / 1000) / 1024)
		metrics.CPUTime[domain.Name] = float64(cpuTime)
	}

	return metrics
}
