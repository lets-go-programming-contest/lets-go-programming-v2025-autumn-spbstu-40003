package main

import (
	"fmt"

	configPkg "github.com/alexpi3/task-8/internal/config"
)

func main() {
	config := configPkg.GetConfig()
	fmt.Print(config.Environment, " ", config.LogLevel)
}
