package config

import (
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
		panic("Failed to open config file: " + err.Error())
	}
	defer file.Close()

	var cfg Settings
	dec := yaml.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		panic("Failed to decode YAML: " + err.Error())
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		panic("Config error: missing input-file or output-file")
	}

	return &cfg, nil
}
