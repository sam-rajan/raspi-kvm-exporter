package config

func setDefaultConfigValue(config *CollectorConfig) {
	if config.Port == nil {
		port := "9000"
		config.Port = &port
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

	if config.Collectors.Raspi.Disk == nil || config.Collectors.Raspi.Disk["enabled"] == nil {
		config.Collectors.Raspi.Disk = map[string]any{
			"enabled": true,
			"devices": config.Collectors.Raspi.Disk["devices"],
		}
	}
}
