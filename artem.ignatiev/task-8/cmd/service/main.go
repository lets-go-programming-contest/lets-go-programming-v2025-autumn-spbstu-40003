package main

import (
	"fmt"
	"log"

	"github.com/kryjkaqq/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
