//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func Load() (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		return nil, errInvalidConfig
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return nil, errInvalidConfig
	}

	return &cfg, nil
}
