package main

import (
	"fmt"

	configPkg "github.com/alexpi3/task-8/internal/config"
)

func main() {
	config := configPkg.GetConfig()
	fmt.Println(config.Environment, config.LogLevel)
}
