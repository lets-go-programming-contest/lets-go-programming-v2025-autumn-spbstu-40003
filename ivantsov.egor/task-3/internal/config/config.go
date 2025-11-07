package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrFileNotSet = errors.New("input-file and output-file must be set")

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic(fmt.Errorf("failed to close config file: %w", cerr))
		}
	}()

	var cfg Config
	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, ErrFileNotSet
	}

	return &cfg, nil
}
