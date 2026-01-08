//go:build dev

package config

import (
	_ "embed"
	"fmt"
)

//go:embed dev.yaml
var devConfig []byte

type devProvider struct{}

func (p *devProvider) GetConfigData() ([]byte, error) {
	if len(devConfig) == 0 {
		return nil, fmt.Errorf("dev config not embedded")
	}

	return devConfig, nil
}

func newDevProvider() provider {
	return &devProvider{}
}
