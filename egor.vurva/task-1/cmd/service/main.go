package main

import (
	"fmt"
	"strconv"
)

func main() {
	var tmpString, operator string

	fmt.Scanln(&tmpString)
	firstValue, err := strconv.Atoi(tmpString)

	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	fmt.Scanln(&tmpString)
	secondValue, err := strconv.Atoi(tmpString)

	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scanln(&operator)

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
