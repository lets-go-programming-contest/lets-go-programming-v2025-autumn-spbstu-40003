package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrInputFileNotSet  = errors.New("input-file must be set")
	ErrOutputFileNotSet = errors.New("output-file must be set")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
		}
	}()

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	if config.InputFile == "" {
		return nil, ErrInputFileNotSet
	}

	if config.OutputFile == "" {
		return nil, ErrOutputFileNotSet
	}

	return &config, nil
}
