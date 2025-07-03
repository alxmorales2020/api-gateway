package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(filePath string) (*GatewayConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config GatewayConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
