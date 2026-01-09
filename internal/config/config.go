package config

type CollectorConfig struct {
	Port       string    `yaml:"port"`
	Collectors Collector `yaml:"collectors"`
}

type Collector struct {
	Kvm   KvmCollectorConfig   `yaml:"kvm"`
	Raspi RaspiCollectorConfig `yaml:"raspi"`
}

type RaspiCollectorConfig struct {
	Enabled bool              `yaml:"enabled"`
	Disk    DiskMetricsConfig `yaml:"disk"`
}

type KvmCollectorConfig struct {
	Enabled bool `yaml:"enabled"`
}

type DiskMetricsConfig struct {
	Enabled bool     `yaml:"enabled"`
	Devices []string `yaml:"devices"`
}
