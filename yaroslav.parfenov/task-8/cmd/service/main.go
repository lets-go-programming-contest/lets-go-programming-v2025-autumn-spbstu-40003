package service

import (
	"fmt"

	"github.com/gituser549/task-8/internal/configprocessing"
)

func main() {
	var cfg configprocessing.BuildConfig

	cfg, err := cfg.GetConfig()
	if err != nil {
		fmt.Println("error getting config: ", err)

		return
	}

	fmt.Println(cfg.Environment, " ", cfg.LogLevel)
}
