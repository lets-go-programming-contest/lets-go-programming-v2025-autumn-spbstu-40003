package main

import (
	"badligyg/task-8/config"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("load config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
