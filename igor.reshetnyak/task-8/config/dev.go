//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var devConfig []byte

type devProvider struct{}

func (p *devProvider) GetConfigData() ([]byte, error) {
	if len(devConfig) == 0 {
		return nil, ErrDevConfigNotEmbedded
	}

	return devConfig, nil
}

func getProvider() provider {
	return &devProvider{}
}
