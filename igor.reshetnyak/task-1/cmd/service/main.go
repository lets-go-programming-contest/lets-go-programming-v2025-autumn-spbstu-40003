package main

import (
	"fmt"
)

func main() {
	var a, b, res int
	var operation string

	_, err := fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("Invalid operetion")
		return
	}

	switch operation {
	case "+":
		res = a + b
	case "-":
		res = a - b
	case "*":
		res = a * b
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		} else {
			res = a / b
		}
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(res)
}
