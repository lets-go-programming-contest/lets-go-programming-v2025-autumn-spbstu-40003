package utils

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errRead  = errors.New("config file read error")
	errParse = errors.New("config file parse error")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errRead
	}

	cfg := Config{"", ""}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errParse
	}

	return &cfg, nil
}
