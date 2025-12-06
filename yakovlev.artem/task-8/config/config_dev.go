//go:build dev
// +build dev

package config

import _ "embed"

//go:embed dev.yaml
var devYAML []byte //nolint:gochecknoglobals // required for embedding

func rawYAML() []byte { return devYAML }
