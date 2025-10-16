package exchangetrade

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func loadConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		panic("error while opening configuration file: " + err.Error())
	}
	defer file.Close()

	var config Config
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		panic("error while parsing configuration file: " + err.Error())
	}

	if config.InputFile == "" {
		panic("error: input-file is not specified in configuration")
	}

	if config.OutputFile == "" {
		panic("error: output-file is not specified in configuration")
	}

	return &config
}
