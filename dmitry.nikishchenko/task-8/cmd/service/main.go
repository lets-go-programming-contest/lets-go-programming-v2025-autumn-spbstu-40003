package main

import (
	"fmt"

	config "github.com/d1mene/task-8/internal/configs"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
