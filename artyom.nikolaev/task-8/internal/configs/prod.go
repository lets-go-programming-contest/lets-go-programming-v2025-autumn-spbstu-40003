//go:build prod

package configs

import _ "embed"

//go:embed config/prod.yaml
var configData []byte
