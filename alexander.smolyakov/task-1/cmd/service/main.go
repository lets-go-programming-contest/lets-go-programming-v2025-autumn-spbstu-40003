package main

import (
	"errors"
	"fmt"
)

func readNumber() (int, error) {
	var n int
	_, err := fmt.Scan(&n)
	return n, err
}

func processNumbers(firstOperand, secondOperand int, operator string) (float64, error) {
	switch operator {
	case "+":
		return float64(firstOperand + secondOperand), nil
	case "-":
		return float64(firstOperand - secondOperand), nil
	case "*":
		return float64(firstOperand * secondOperand), nil
	case "/":
		if secondOperand == 0 {
			return 0, errors.New("Division by zero")
		}
		return float64(firstOperand) / float64(secondOperand), nil
	default:
		return 0, errors.New("Invalid operation")
	}
}


func main() {
	var (
		firstOperand, secondOperand int
		result float64
		operator string
	)

	firstOperand, err := readNumber()
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	secondOperand, err = readNumber()
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Error reading operator")
		return
	}

	result, err = processNumbers(firstOperand, secondOperand, operator)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

