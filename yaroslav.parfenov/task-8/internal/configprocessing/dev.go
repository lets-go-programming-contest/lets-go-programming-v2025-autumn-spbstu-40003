//go:build dev

package configprocessing

import _ "embed"

//go:embed dev.yaml
var ConfigFileContent []byte
