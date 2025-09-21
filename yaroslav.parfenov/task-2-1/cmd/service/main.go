package main

import (
	"fmt"

	"github.com/gituser549/task-2-1/internal/solution"
)

func main() {
	var (
		numDepartments int
	)

	_, err := fmt.Scanln(&numDepartments)
	if err != nil {
		fmt.Println("Invalid number of departments")
	}

	err = solution.ProcessEmployees(&numDepartments)
	if err != nil {
		fmt.Println(err.Error())
	}
}
