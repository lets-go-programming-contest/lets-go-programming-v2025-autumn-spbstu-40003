package config

// Файл с парсингом конфигурации из .yaml

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

func InitConfig() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config .yaml file")
	flag.Parse()

	return configPath
}

func ParseConfig() (string, string) {
	file, err := os.Open(InitConfig())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := Config{}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config.Input, config.Output
}
