package main

import (
	"fmt"
	"strconv"
)

func main() {
	var tmpString, operator string

	_, err := fmt.Scanln(&tmpString)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	firstValue, err := strconv.Atoi(tmpString)

	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&tmpString)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	secondValue, err := strconv.Atoi(tmpString)

	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operator)
	if err != nil {
		fmt.Println("Invalid first operand")
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
			return
		} else {
			fmt.Println(firstValue / secondValue)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
