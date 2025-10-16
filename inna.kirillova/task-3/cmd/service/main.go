package main

import (
"flag"

"github.com/kirinnah/task-3/internal/exchangetrade"
)

func main() {
configPath := flag.String("config", "config.yaml", "path to configuration file")
flag.Parse()

exchangetrade.Process(*configPath)
}
