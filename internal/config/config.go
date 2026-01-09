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
	Enabled *bool          `yaml:"enabled,omitempty"`
	Disk    map[string]any `yaml:"disk,omitempty"`
}

type KvmCollectorConfig struct {
	Enabled *bool `yaml:"enabled,omitempty"`
}
