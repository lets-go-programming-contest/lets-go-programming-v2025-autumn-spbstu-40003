package configprocessing

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type BuildConfig struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func (bc *BuildConfig) GetConfig() (BuildConfig, error) {
	var cfg BuildConfig

	if err := yaml.Unmarshal(ConfigFileContent, &cfg); err != nil {
		return BuildConfig{}, fmt.Errorf("error parsing config file: %w", err)
	}

	return cfg, nil
}
