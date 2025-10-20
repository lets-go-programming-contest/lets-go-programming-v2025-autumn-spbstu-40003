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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	if config.InputFile == "" {
		return nil, fmt.Errorf("input-file must be set")
	}

	if config.OutputFile == "" {
		return nil, fmt.Errorf("output-file must be set")
	}

	return &config, nil
}
