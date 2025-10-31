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

func Read(path string) (*Settings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("close error:", cerr)
		}
	}()

	var cfg Settings
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, fmt.Errorf("config error: missing input-file or output-file")
	}

	return &cfg, nil
}
