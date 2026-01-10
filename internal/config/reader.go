package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func (config *CollectorConfig) loadConfig(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
