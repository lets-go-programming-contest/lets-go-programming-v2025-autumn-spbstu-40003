//go:build dev

package config

import _ "embed"

//go:embed configs/dev.yaml
var ConfigData []byte
