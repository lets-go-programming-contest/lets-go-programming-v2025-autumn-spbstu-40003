package configreader

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var errEmptyFlag = errors.New("config flag is empty")

type Config struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func Parse() (Config, error) {
	var config Config

	configPath := flag.String("config", "", "input path config file")
	flag.Parse()

	if *configPath == "" {
		return config, errEmptyFlag
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return config, fmt.Errorf("failed reading file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed deserialization YAML: %w", err)
	}

	return config, nil
}
