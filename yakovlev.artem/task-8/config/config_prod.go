//go:build !dev
// +build !dev

package config

import _ "embed"

var prodYAML []byte

func rawYAML() []byte { return prodYAML }
