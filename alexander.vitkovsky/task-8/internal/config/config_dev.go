//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var rawConfig []byte

func load() Config {
	var config Config

	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		panic(err)
	}

	return config
}
