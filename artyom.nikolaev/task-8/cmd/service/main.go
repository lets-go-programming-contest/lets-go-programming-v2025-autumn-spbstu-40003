package main

import (
	"fmt"

	"github.com/ArtttNik/task-8/internal/configs"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
