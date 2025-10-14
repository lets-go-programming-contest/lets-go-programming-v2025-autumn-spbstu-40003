package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open config: %w", err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(fmt.Errorf("failed to close config: %w", err))
		}
	}()

	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return &cfg, fmt.Errorf("failed to decode YAML: %w", err)
	}
	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return &cfg, fmt.Errorf("input-file and output-file must be set")
	}
	return &cfg, nil
}
