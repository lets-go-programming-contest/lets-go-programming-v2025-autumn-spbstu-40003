//go:build !dev
// +build !dev

package config

import _ "embed"

//nolint:gochecknoglobals // required by go:embed
//go:embed prod.yaml
var prodYAML []byte

func rawYAML() []byte { return prodYAML }
