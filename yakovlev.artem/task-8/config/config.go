package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(rawYAML(), &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Environment == "" || cfg.LogLevel == "" {
		return Config{}, fmt.Errorf("bad config: environment/log_level is empty")
	}
	return cfg, nil
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
