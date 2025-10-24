package main

import "fmt"

const (
	minTemperature   = 15
	maxTemperature   = 30
	invalidIndicator = -1
)

func main() {
	var depCount int

	if _, scanError := fmt.Scan(&depCount); scanError != nil {
		fmt.Printf("Error reading department count: %v\n", scanError)
		fmt.Println(invalidIndicator)

		return
	}

	for range depCount {
		var empCount int

		if _, scanError := fmt.Scan(&empCount); scanError != nil {
			fmt.Println(invalidIndicator)

			return
		}

		processDepartment(empCount)
	}
}

func processDepartment(empCount int) {
	currentMin := minTemperature
	currentMax := maxTemperature
	isPossible := true

	for range empCount {
		var condition string
		var desiredTemp int

		if _, err := fmt.Scan(&condition, &desiredTemp); err != nil {
			fmt.Println(invalidIndicator)

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
			fmt.Println(invalidIndicator)
			isPossible = false
			continue
		}

		if currentMin <= currentMax {
			fmt.Println(currentMin)
		} else {
			fmt.Println(invalidIndicator)
			isPossible = false
		}
	}
}
