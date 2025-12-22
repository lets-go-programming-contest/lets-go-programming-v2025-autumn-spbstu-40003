package configs

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed yaml_configs/*.yaml
var yamlFiles embed.FS

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	data, err := yamlFiles.ReadFile("yaml_configs/" + configFileName())
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &cfg, nil
}
