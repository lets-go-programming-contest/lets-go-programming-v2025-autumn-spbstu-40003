//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var data []byte

func Load() (*Config, error) {
	return parseConfig(data)
}
