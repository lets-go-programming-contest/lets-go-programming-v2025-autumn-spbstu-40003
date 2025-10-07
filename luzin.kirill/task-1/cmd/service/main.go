package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var (
		input, operation string
		a, b             int
		result           float64
	)

	_, err := fmt.Scan(&input)

	if err != nil {
		fmt.Println("Invalid first operand")
		os.Exit(0)
	}

	a, err = strconv.Atoi(input)

	if err != nil {
		fmt.Println("Invalid first operand")
		os.Exit(0)
	}

	_, err = fmt.Scan(&input)

	if err != nil {
		fmt.Println("Invalid second operand")
		os.Exit(0)
	}

	b, err = strconv.Atoi(input)

	if err != nil {
		fmt.Println("Invalid second operand")
		os.Exit(0)
	}

	_, err = fmt.Scan(&operation)

	if err != nil {
		fmt.Println("Invalid operation ")
		os.Exit(0)
	}

	switch operation {
	case "+":
		result = float64(a + b)

	case "-":
		result = float64(a - b)

	case "*":
		result = float64(a * b)

	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			os.Exit(0)
		}
		result = float64(a) / float64(b)

	default:
		fmt.Println("Invalid operation")
		os.Exit(0)
	}

	fmt.Println(result)
}
