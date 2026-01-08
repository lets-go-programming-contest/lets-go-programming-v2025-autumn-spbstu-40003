//go:build dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func Load() (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return nil, errInvalidConfig
	}

	return &cfg, nil
}
