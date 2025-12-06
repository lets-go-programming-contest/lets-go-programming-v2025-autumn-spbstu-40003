//go:build prod || !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var configFile []byte

func Load() (Config, error) {
	return parseYAML(configFile)
}
