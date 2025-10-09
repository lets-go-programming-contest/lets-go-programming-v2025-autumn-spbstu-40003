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
	data, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}

	defer func(data *os.File) {
		err := data.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file: %w", err))
		}
	}(data)

	decoder := yaml.NewDecoder(data)

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		return &config, fmt.Errorf("failed to decode file: %w", err)
	}

	if config.InputFile == "" || config.OutputFile == "" {
		return &config, fmt.Errorf("input-file and output-file must be set: %w", err)
	}

	return &config, nil
}
