package main

import (
	"fmt"

	temperature "github.com/d1mene/task-2-1/internal"
)

func main() {
	err := temperature.TemperatureControl()
	if err != nil {
		fmt.Println(err)
	}
}
