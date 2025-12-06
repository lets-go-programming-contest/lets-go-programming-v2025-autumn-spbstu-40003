//go:build !dev
// +build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodYAML []byte //nolint:gochecknoglobals // required for embedding

func rawYAML() []byte { return prodYAML }
