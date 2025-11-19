package main

import (
	"errors"
	"fmt"
)

var (
	ErrDepart   = errors.New("departments error")
	ErrEmployee = errors.New("employee error")
	ErrTemp     = errors.New("incorrect temperature")
	ErrSymbol   = errors.New("incorrect symbol")
)

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil || departments < 1 || departments > 1000 {
		fmt.Println(-1)
		return
	}

	for i := 0; i < departments; i++ {
		var employee int

		if _, err := fmt.Scan(&employee); err != nil || employee < 1 || employee > 1000 {
			fmt.Println(-1)
			return
		}

		departmentOptimalTemp(employee)
	}
}

func departmentOptimalTemp(employee int) {
	minTemp := 15
	maxTemp := 30
	inRangeTemp := true

	for i := 0; i < employee; i++ {
		var symbol string
		var newTemp int

		if _, err := fmt.Scan(&symbol); err != nil {
			fmt.Println(-1)
			return
		}

		if _, err := fmt.Scan(&newTemp); err != nil {
			fmt.Println(-1)
			return
		}

		if !inRangeTemp {
			fmt.Println(-1)
			continue
		}

		switch symbol {
		case ">=":
			if newTemp > minTemp {
				minTemp = newTemp
			}
		case "<=":
			if newTemp < maxTemp {
				maxTemp = newTemp
			}
		default:
			fmt.Println(-1)
			return
		}

		if minTemp <= maxTemp {
			fmt.Println(minTemp)
		} else {
			fmt.Println(-1)
			inRangeTemp = false
		}
	}
}
