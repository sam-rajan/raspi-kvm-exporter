package config

func setDefaultConfigValue(config *CollectorConfig) {
	if config.Port == "" {
		config.Port = "9000"
	}

	if config.Collectors.Kvm.Enabled == nil {
		config.Collectors.Kvm.Enabled = new(bool)
		*config.Collectors.Kvm.Enabled = true
	}

	if config.Collectors.Raspi.Enabled == nil {
		config.Collectors.Raspi.Enabled = new(bool)
		*config.Collectors.Raspi.Enabled = true
	}

	if config.Collectors.Raspi.Disk == nil {
		config.Collectors.Raspi.Enabled = new(bool)
		*config.Collectors.Raspi.Enabled = true
	}

	if config.Collectors.Raspi.Disk == nil || config.Collectors.Raspi.Disk["Enabled"] == nil {
		config.Collectors.Raspi.Disk = map[string]any{
			"Enabled": true,
			"Devices": config.Collectors.Raspi.Disk["Devices"],
		}
	}
}
