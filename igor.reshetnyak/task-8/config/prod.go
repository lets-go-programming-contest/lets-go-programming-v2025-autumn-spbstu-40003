//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodConfig []byte

type prodProvider struct{}

func (p *prodProvider) GetConfigData() ([]byte, error) {
	if len(prodConfig) == 0 {
		return nil, ErrProdConfigNotEmbedded
	}

	return prodConfig, nil
}

func getProvider() provider {
	return &prodProvider{}
}
