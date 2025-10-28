package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputPath  string `yaml:"input-file"`
	OutputPath string `yaml:"output-file"`
}

var (
	ErrNoInputFile  = errors.New("input-file not set")
	ErrNoOutputFile = errors.New("output-file not set")
)

func Load(path *string) (*Config, error) {
	content, err := os.ReadFile(*path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("cannot parse YAML: %w", err)
	}

	if cfg.InputPath == "" {
		return nil, ErrNoInputFile
	}
	if cfg.OutputPath == "" {
		return nil, ErrNoOutputFile
	}

	return cfg, nil
}
