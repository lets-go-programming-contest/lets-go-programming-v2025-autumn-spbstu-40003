package main

import "fmt"

const (
	minTemperature   = 15
	maxTemperature   = 30
	invalidIndicator = -1
)

func main() {
	var elementCount int

	if _, err := fmt.Scan(&elementCount); err != nil {
		fmt.Printf("Error reading input (context: elementCount): %v\n", err)

		return
	}

	for i := 1; i <= elementCount; i++ {
		var employeeCount int

		if _, err := fmt.Scan(&employeeCount); err != nil {
			fmt.Printf("Error reading input (context: employeeCount in department %d): %v\n", i, err)

			return
		}

		processDepartment(i, employeeCount)
	}
}

func processDepartment(departmentIndex, employeeCount int) {
	currentMin := minTemperature
	currentMax := maxTemperature
	isPossible := true

	for j := 1; j <= employeeCount; j++ {
		var condition string
		var desiredTemp int

		if _, err := fmt.Scan(&condition, &desiredTemp); err != nil {
			fmt.Printf("Error reading input (context: condition or desiredTemp for employee %d in department %d): %v\n", j, departmentIndex, err)

			return
		}

		if !isPossible {
			fmt.Println(invalidIndicator)
			continue
		}

		switch condition {
		case ">=":
			if desiredTemp > currentMin {
				currentMin = desiredTemp
			}
		case "<=":
			if desiredTemp < currentMax {
				currentMax = desiredTemp
			}
		default:
			fmt.Printf("Error: invalid condition '%s' (employee %d, department %d)\n", condition, j, departmentIndex)
			fmt.Println(invalidIndicator)
			isPossible = false
			continue
		}

		if currentMin <= currentMax {
			fmt.Println(currentMin)
		} else {
			fmt.Printf("Error: impossible temperature range (min=%d, max=%d) in department %d\n",
				currentMin, currentMax, departmentIndex)
			fmt.Println(invalidIndicator)
			isPossible = false
		}
	}
}
