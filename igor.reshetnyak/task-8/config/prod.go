//go:build !dev

package config

import (
	_ "embed"
	"fmt"
)

//go:embed prod.yaml
var prodConfig []byte

type prodProvider struct{}

func (p *prodProvider) GetConfigData() ([]byte, error) {
	if len(prodConfig) == 0 {
		return nil, fmt.Errorf("prod config not embedded")
	}

	return prodConfig, nil
}

func newProdProvider() provider {
	return &prodProvider{}
}
