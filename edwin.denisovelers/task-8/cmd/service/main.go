package main

import (
	"fmt"

	"github.com/wedwincode/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("environment: %s, log_level: %s", cfg.Environment, cfg.LogLevel)
}
