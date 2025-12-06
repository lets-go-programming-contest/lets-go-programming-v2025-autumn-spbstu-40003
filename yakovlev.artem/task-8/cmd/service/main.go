package main

import (
	"fmt"

	"github.com/nxgmvw/task-8/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
