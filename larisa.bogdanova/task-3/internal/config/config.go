package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadSettings(path string) (*Settings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	var settings Settings
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&settings); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	if settings.InputFile == "" || settings.OutputFile == "" {
		return nil, fmt.Errorf("input-file and output-file must be set")
	}

	return &settings, nil
}
