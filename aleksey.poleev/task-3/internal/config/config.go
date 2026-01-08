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
	data, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if cerr := data.Close(); cerr != nil {
			panic(fmt.Errorf("failed to close file: %w", cerr))
		}
	}()

	decoder := yaml.NewDecoder(data)

	var cfg Config

	if err := decoder.Decode(&cfg); err != nil {
		return &cfg, fmt.Errorf("failed to decode file: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, ErrFileNotSet
	}

	return &cfg, nil
}
