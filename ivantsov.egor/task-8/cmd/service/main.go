package main

import (
	"fmt"
	"log"

	"github.com/lets-go-programming-contest/lets-go-programming-v2025-autumn-spbstu-40003/ivantsov.egor/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
