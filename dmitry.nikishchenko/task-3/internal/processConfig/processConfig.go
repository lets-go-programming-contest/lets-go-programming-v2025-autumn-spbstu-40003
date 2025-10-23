package processconfig

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	errInput  = errors.New("invalid config: missing input-file")
	errOutput = errors.New("invalid config: missing output-file")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig() (*Config, error) {
	configPath := flag.String("config", "config.yaml", "path to .yaml config file")
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("cannot close file: %v", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot raed config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	if cfg.InputFile == "" {
		return nil, errInput
	}

	if cfg.OutputFile == "" {
		return nil, errOutput
	}

	return &cfg, nil
}
