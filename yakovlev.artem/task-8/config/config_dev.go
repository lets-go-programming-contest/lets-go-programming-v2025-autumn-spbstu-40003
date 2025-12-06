//go:build dev
// +build dev

package config

import _ "embed"

//nolint:gochecknoglobals // required by go:embed
//go:embed dev.yaml
var devYAML []byte

func rawYAML() []byte { return devYAML }
