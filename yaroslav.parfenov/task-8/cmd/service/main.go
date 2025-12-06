package main

import (
	"fmt"

	"github.com/gituser549/task-8/internal/configprocessing"
)

func main() {
	var cfg configprocessing.BuildConfig

	cfg, err := cfg.GetConfig()
	if err != nil {
		fmt.Printf("error getting config: %v\n", err)

		return
	}

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
