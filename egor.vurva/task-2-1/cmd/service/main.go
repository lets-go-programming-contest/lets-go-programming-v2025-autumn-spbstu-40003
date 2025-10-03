package main

import (
	"fmt"

	"github.com/Vurvaa/task-2-1/internal/temperature"
)

func main() {
	var deptCount, employeeCount int

	_, err := fmt.Scan(&deptCount)
	if err != nil {
		return
	}

	for deptCount > 0 {
		_, err = fmt.Scan(&employeeCount)
		if err != nil {
			return
		}
		temperature.CheckRange(employeeCount)
		deptCount--
	}
}
