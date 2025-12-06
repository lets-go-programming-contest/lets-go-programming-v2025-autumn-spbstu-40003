package config

import (
	"fmt"
	"strings"
)

type Config struct {
	Environment string
	LogLevel    string
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (Config, error) {
	return parseYAML(rawYAML())
}

func parseYAML(b []byte) (Config, error) {
	m := map[string]string{}

	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		i := strings.Index(line, ":")
		if i < 0 {
			continue
		}

		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+1:])
		val = strings.Trim(val, `"'`)

		m[key] = val
	}

	cfg := Config{
		Environment: m["environment"],
		LogLevel:    m["log_level"],
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return Config{}, fmt.Errorf("bad config: expected keys environment and log_level")
	}

	return cfg, nil
}
