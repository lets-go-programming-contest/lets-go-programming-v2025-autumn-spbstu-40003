//go:build prod || !dev

package configs

func init() {
	configFileName = "prod.yaml"
}
