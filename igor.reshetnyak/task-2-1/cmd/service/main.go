package main

import (
	"errors"
	"fmt"
)

var (
	ErrDepart   = errors.New("deparments error")
	ErrEmployee = errors.New("employee error")
	ErrTemp     = errors.New("incorrect temperature")
	ErrSymbol   = errors.New("incorrerc symbol")
)

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil || departments < 1 || departments > 1000 {
		fmt.Println(ErrDepart, err)
		return
	}

	for range departments {
		var employee int

		if _, err := fmt.Scan(&employee); err != nil || employee < 1 || employee > 1000 {
			fmt.Println(ErrEmployee, err)
			return
		}

		deparmentOptimalTemp(employee)
	}
}

func deparmentOptimalTemp(employee int) {
	minTemp := 15
	maxTemp := 30
	inRangeTemp := true

	for range employee {
		var symbol string
		var newTemp int

		if _, err := fmt.Scan(&symbol); err != nil || (symbol != "<=" && symbol != ">=") {
			fmt.Println(ErrSymbol, err)
			return
		}

		if _, err := fmt.Scan(&newTemp); err != nil {
			fmt.Println(ErrTemp, err)
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
