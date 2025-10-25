package config

import (
	"os"

	"github.com/stretchr/testify/assert/yaml"
)

type conf struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetConfig(configPATH string) (conf, error) {
	str, err := os.ReadFile(configPATH)
	if err != nil {
		panic("no such file or directory")
	}

	var cfg conf
	if err := yaml.Unmarshal(str, &cfg); err != nil {
		panic("did not find expected key")
	}

	return cfg, nil
}
