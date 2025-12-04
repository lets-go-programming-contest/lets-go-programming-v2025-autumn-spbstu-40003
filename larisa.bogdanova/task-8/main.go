package main

import (
	"fmt"

	"badligyg/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config error: %w", err)

		return
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
