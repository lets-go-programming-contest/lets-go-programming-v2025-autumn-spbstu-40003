//go:build !prod

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var configFile []byte

func Load() (Config, error) {
	return parseYAML(configFile)
}
