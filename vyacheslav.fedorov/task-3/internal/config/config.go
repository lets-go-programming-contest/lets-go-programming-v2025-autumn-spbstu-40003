package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ParseConfiguration(path string) (*Configuration, error) {
	data, err := os.ReadFile(path)
	if err != nil {

		return nil, fmt.Errorf("error reading config file %q: %w", path, err)
	}

	var cfg Configuration

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {

		return nil, fmt.Errorf("error parsing config file %q: %w", path, err)

	}

	return &cfg, nil

}
