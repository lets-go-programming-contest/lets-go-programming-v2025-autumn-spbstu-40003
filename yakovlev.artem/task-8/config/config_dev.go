//go:build dev
// +build dev

package config

import _ "embed"

var devYAML []byte

func rawYAML() []byte { return devYAML }
