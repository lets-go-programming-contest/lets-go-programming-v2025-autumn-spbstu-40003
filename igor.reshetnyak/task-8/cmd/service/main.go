package main

import (
	"fmt"
	"os"

	"github.com/ReshetnyakIgor/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
