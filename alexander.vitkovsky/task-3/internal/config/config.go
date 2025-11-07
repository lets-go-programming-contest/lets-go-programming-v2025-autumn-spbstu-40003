package config

import (
	"flag"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func ParseConfig() (string, string) {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config .yaml file")
	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		panic("failed to open input file" + err.Error())
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic("failed to read config file" + err.Error())
	}

	if err := file.Close(); err != nil {
		panic("failed to close config file" + err.Error())
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		panic("failed to parse config file" + err.Error())
	}

	return config.Input, config.Output
}
