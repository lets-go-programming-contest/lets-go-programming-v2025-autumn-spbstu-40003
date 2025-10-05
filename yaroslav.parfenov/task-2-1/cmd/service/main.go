package main

import (
	"fmt"

	"github.com/gituser549/task-2-1/internal/temperatureoptimizer"
)

func main() {
	var numDepartments int

	_, err := fmt.Scanln(&numDepartments)
	if err != nil {
		fmt.Println("invalid number of departments")
	}

	err = temperatureoptimizer.ProcessDepartments(&numDepartments)
	if err != nil {
		fmt.Println(err.Error())
	}
}
