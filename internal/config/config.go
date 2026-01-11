package config

import "log"

type CollectorConfig struct {
	Port       string    `yaml:"port"`
	Collectors Collector `yaml:"collectors"`
}

type Collector struct {
	Kvm   KvmCollectorConfig   `yaml:"kvm"`
	Raspi RaspiCollectorConfig `yaml:"raspi"`
}

type RaspiCollectorConfig struct {
	Enabled *bool          `yaml:"enabled,omitempty"`
	Disk    map[string]any `yaml:"disk,omitempty"`
}

type KvmCollectorConfig struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}

func NewCollectorConfig(configFile *string) *CollectorConfig {
	collectorConfig := &CollectorConfig{
		Port: "9000",
		Collectors: Collector{
			Kvm:   KvmCollectorConfig{},
			Raspi: RaspiCollectorConfig{},
		},
	}

	if configFile != nil {
		err := collectorConfig.loadConfig(*configFile)
		if err != nil {
			log.Println(err.Error())
		}
	}

	setDefaultConfigValue(collectorConfig)
	return collectorConfig
}
