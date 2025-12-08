//go:build !dev

package configprocessing

import _ "embed"

//go:embed prod.yaml
var ConfigFileContent []byte
