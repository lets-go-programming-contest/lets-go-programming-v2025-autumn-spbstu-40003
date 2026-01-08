package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

type provider interface {
	GetConfigData() ([]byte, error)
}

var currentProvider provider

func Load() (*Config, error) {
	if currentProvider == nil {
		return nil, fmt.Errorf("config provider not initialized")
	}

	data, err := currentProvider.GetConfigData()
	if err != nil {
		return nil, fmt.Errorf("failed to get config data: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Environment == "" {
		return fmt.Errorf("environment is required")
	}
	if c.LogLevel == "" {
		return fmt.Errorf("log_level is required")
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("%s %s", c.Environment, c.LogLevel)
}
