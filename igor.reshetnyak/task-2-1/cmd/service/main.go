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

const (
	minTempConst = 15
	maxTempConst = 30
	inRangeConst = true
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

		departmentOptimalTemp(employee)
	}
}

func departmentOptimalTemp(employee int) {
	minTemp := minTempConst
	maxTemp := maxTempConst
	inRangeTemp := inRangeConst

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

		updateAndCheck(symbol, newTemp, &minTemp, &maxTemp, &inRangeTemp)
	}
}

func updateAndCheck(symbol string, newTemp int, minTemp *int, maxTemp *int, inRangeTemp *bool) {
	switch symbol {
	case ">=":
		if newTemp >= *minTemp {
			*minTemp = newTemp
		}
	case "<=":
		if newTemp <= *maxTemp {
			*maxTemp = newTemp
		}
	default:
		fmt.Println(-1)

		*inRangeTemp = false

		return
	}

	if *minTemp <= *maxTemp {
		fmt.Println(*minTemp)
	} else {
		fmt.Println(-1)

		*inRangeTemp = false
	}
}
