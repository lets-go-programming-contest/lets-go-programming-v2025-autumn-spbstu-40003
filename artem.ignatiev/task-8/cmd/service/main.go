package main

import (
	"fmt"
	"log"

	configs "github.com/kryjkaqq/task-8/internal/config"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
