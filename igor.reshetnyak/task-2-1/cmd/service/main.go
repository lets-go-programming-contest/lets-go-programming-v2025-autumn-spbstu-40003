package main

import (
	"errors"
	"fmt"
)

var (
	ErrDepart   = errors.New("departments error")
	ErrEmployee = errors.New("employee error")
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
	tempState := &TempState{
		Min:     15,
		Max:     30,
		InRange: true,
	}

	for range employee {
		if !tempState.ProcessEmployee() {
			return
		}
	}
}

type TempState struct {
	Min     int
	Max     int
	InRange bool
}

func (tempState *TempState) ProcessEmployee() bool {
	if !tempState.InRange {
		fmt.Println(-1)
		return true
	}

	symbol, newTemp := tempState.readInput()
	if symbol == "" {
		return false
	}

	tempState.updateTemperatures(symbol, newTemp)
	tempState.checkRange()

	return true
}

func (tempState *TempState) readInput() (string, int) {
	var symbol string
	var newTemp int

	if _, err := fmt.Scan(&symbol); err != nil || (symbol != "<=" && symbol != ">=") {
		fmt.Println(-1)
		return "", 0
	}

	if _, err := fmt.Scan(&newTemp); err != nil {
		fmt.Println(-1)
		return "", 0
	}

	return symbol, newTemp
}

func (tempState *TempState) updateTemperatures(symbol string, newTemp int) {
	switch symbol {
	case ">=":
		if newTemp > tempState.Min {
			tempState.Min = newTemp
		}
	case "<=":
		if newTemp < tempState.Max {
			tempState.Max = newTemp
		}
	default:
		fmt.Println(-1)
		tempState.InRange = false
	}
}

func (tempState *TempState) checkRange() {
	if tempState.Min <= tempState.Max {
		fmt.Println(tempState.Min)
	} else {
		fmt.Println(-1)
		tempState.InRange = false
	}
}
