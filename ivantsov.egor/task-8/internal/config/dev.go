//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var data []byte

func Load() (*Config, error) {
	return parseConfig(data)
}
