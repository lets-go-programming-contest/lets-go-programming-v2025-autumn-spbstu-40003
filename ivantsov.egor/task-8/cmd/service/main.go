package main

import (
	"fmt"
	"log"

<<<<<<< HEAD
	"github.com/lets-go-programming-contest/lets-go-programming-v2025-autumn-spbstu-40003/ivantsov.egor/task-8/internal/config"
=======
	"github.com/Egor1726/lets-go-programming-contest/ivantsov.egor/task-8/internal/config"
>>>>>>> 86a1198 ([TASK-8] Initial config structure)
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
