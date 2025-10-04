package main

import (
	"fmt"

	"github.com/Vurvaa/task-2-1/internal/temperature"
)

func main() {
	var deptCount, employeeCount int

	_, err := fmt.Scan(&deptCount)
	if err != nil {
		fmt.Println("failed to scan deptCount:", err)

		return
	}

	for deptCount > 0 {
		_, err = fmt.Scan(&employeeCount)
		if err != nil {
			fmt.Println("failed to scan employeeCount:", err)

			return
		}

		err = temperature.CheckRange(employeeCount)
		if err != nil {
			fmt.Println("check range failed:", err)

			return
		}

		deptCount--
	}
}
