package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

func (config *CollectorConfig) LoadConfig(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, config)
}
