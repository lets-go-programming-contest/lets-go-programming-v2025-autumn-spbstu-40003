package main

import (
	"fmt"
	"os"

	"github.com/nxgmvw/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
