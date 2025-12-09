package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var (
	ErrConfigRequiredField = errors.New("config: required field is empty")
)

func Load() (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal yaml: %w", err)
	}

	if cfg.Environment == "" {
		return Config{}, fmt.Errorf("%w: environment", ErrConfigRequiredField)
	}
	if cfg.LogLevel == "" {
		return Config{}, fmt.Errorf("%w: log_level", ErrConfigRequiredField)
	}

	return cfg, nil
}
