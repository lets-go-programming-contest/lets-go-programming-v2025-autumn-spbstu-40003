package processconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetConfig(cfgPath string) (Config, error) {
	var cfg Config

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return cfg, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		err := cfgFile.Close()
		if err != nil {
			panic(fmt.Errorf("error closing config file: %w", err))
		}
	}()

	yamlDecoder := yaml.NewDecoder(cfgFile)

	err = yamlDecoder.Decode(&cfg)

	if err != nil {
		return cfg, fmt.Errorf("error decoding config file: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return cfg, fmt.Errorf("invalid config file: missing input-file or output-dir: %w", err)
	}

	return cfg, nil
}
