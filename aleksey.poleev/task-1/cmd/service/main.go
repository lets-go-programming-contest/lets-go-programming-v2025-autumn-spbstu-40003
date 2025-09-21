package main

import (
	"fmt"
	"strconv"
)

func main() {
	var in, op string

	fmt.Scan(&in)
	a, err := strconv.Atoi(in)
	if err != nil {
		fmt.Println("Invalid first operand")
	}
	fmt.Scan(&in)
	b, err := strconv.Atoi(in)
	if err != nil {
		fmt.Println("Invalid second operand")
	}
	fmt.Scan(&op)

	switch {
	case op == "+":
		fmt.Println(a + b)
	case op == "-":
		fmt.Println(a - b)
	case op == "*":
		fmt.Println(a * b)
	case op == "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
