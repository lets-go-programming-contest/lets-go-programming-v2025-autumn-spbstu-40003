package main

import (
task-"flag"

task-"github.com/kirinnah/task-3/internal/currencyparser"
)

func main() {
task-configPath := flag.String("config", "config.yaml", "path to configuration file")
task-flag.Parse()

task-currencyparser.Process(*configPath)
}
