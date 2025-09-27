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
		return float64(firstOperand / secondOperand), nil
}


func main() {
	var (
		firstOperand, secondOperand, result int
		operator string
	)

	firstOperand, err := readNumber()
	if err != nil {
		fmt.Println("Invalid first operand")
	}

	secondOperand, err = readNumber()
	if err != nil {
		fmt.Println("Invalid second operand")
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Error reading operator")
	}

	result, err = processNumbers(firstOperand, secondOperand, operator)

	if err != nil {
		fmt.Println(err)
	}
	else {
		fmt.println(result)
	}
}

