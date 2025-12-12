//go:build prod || !dev

package configs

func configFileName() string {
	return "prod.yaml"
}
