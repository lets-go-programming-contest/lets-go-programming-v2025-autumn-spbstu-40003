package main

import "fmt"

const (
	minTemp      = 15
	maxTemp      = 30
	invalidValue = -1
)

func main() {
	var numOfDep int

	_, err := fmt.Scan(&numOfDep)
	if err != nil {
		fmt.Println(invalidValue)
		return
	}

	for range numOfDep {
		err := handleDepartment()
		if err != nil {
			fmt.Println(invalidValue)
			return
		}
	}
}

func handleDepartment() error {
	var numOfEmpl int

	_, err := fmt.Scan(&numOfEmpl)
	if err != nil {
		return fmt.Errorf("Error while reading departments: %w", err)
	}

	lowerLimit, upperLimit := minTemp, maxTemp
	flag := true

	for range numOfEmpl {
		var (
			sign string
			temp int
		)

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			return fmt.Errorf("Error while reading temperature preference: %w", err)
		}

		if !flag {
			fmt.Println(invalidValue)

			continue
		}

		if sign == ">=" {
			if temp > lowerLimit {
				lowerLimit = temp
			}
		} else {
			if temp < upperLimit {
				upperLimit = temp
			}
		}

		if lowerLimit <= upperLimit {
			fmt.Println(lowerLimit)
		} else {
			fmt.Println(invalidValue)

			flag = false
		}
	}

	return nil
}
