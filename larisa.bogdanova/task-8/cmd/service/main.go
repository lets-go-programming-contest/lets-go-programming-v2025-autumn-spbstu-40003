package main

import (
	"fmt"
	"os"

	"larisa.bogdanova/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("load config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
