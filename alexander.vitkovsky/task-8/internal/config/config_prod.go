//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var rawConfig []byte

func load() Config {
	var config Config

	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		panic(err)
	}

	return config
}
