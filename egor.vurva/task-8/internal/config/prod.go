//go:build prod || !dev

package config

import _ "embed"

//go:embed prod.yaml
var content []byte
