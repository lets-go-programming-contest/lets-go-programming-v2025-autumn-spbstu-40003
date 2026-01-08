package main

import (
	"fmt"
	"log"

	"github.com/Egor1726/lets-go-programming-contest/ivantsov.egor/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
