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

func getProvider() provider {
	return nil
}

func Load() (*Config, error) {
	p := getProvider()
	if p == nil {
		return nil, ErrProviderNotInitialized
	}

	data, err := p.GetConfigData()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {

	if c.Environment == "" {
		return ErrEnvironmentRequired
	}

	if c.LogLevel == "" {
		return ErrLogLevelRequired
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("%s %s", c.Environment, c.LogLevel)
}
