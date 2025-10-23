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

type Settings struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadSettings(path string) (*Settings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
		}
	}()

	var settings Settings

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&settings); err != nil {
		return nil, fmt.Errorf("decode config file: %w", err)
	}

	if settings.InputFile == "" {
		return nil, ErrInputFileNotSet
	}

	if settings.OutputFile == "" {
		return nil, ErrOutputFileNotSet
	}

	return &settings, nil
}
