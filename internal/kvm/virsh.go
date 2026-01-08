package kvm

import (
	"log"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
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

	l := libvirt.NewWithDialer(dialers.NewLocal())
	if err := l.Connect(); err != nil {
		log.Printf("Failed to connect: %s", err.Error())
		return metrics
	}

	domains, _, err := l.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive)
	if err != nil {
		log.Printf("failed to retrieve domains: %v", err)
		return metrics
	}

	for _, domain := range domains {
		state, _, mem, _, cpuTime, _ := l.DomainGetInfo(domain)
		metrics.VMStatus[domain.Name] = float64(state)
		memKb := (mem / 1000) + (mem % 1000)
		memMb := (memKb / 1024) + (memKb % 1024)
		metrics.MemoryStatus[domain.Name] = float64(memMb)
		metrics.CPUTime[domain.Name] = float64(cpuTime)
	}

	return metrics
}
