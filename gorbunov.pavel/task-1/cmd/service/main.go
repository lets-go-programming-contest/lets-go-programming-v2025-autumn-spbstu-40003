package main

import (
	"fmt"
	"strconv"
)

func main() {
	var inputString, operator string

	_, err := fmt.Scanln(&inputString)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	firstValue, err := strconv.Atoi(inputString)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&inputString)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	secondValue, err := strconv.Atoi(inputString)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operator)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(firstValue + secondValue)
	case "-":
		fmt.Println(firstValue - secondValue)
	case "*":
		fmt.Println(firstValue * secondValue)
	case "/":
		if secondValue == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(firstValue / secondValue)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
