package main

import (
	"fmt"

	"github.com/Vurvaa/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("invalid state before load config: %s\n", err)

		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
