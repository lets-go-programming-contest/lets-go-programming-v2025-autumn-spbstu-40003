package main

import (
	"fmt"
	"os"

	"github.com/ReshetnyakIgor/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cfg.String())
}
