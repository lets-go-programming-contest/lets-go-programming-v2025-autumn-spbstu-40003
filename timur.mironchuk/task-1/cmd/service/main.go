package main

import (
	"fmt"
	"os"
)

func main() {
	var input1 int
	var input2 int
	var operator string

	fmt.Println("Введите первое число: ")
	if _, err := fmt.Scan(&input1); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	fmt.Println("Введите второе число: ")
	if _, err := fmt.Scan(&input2); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Println("Введите оператор: ")
	if _, err := fmt.Scan(&operator); err != nil {
		fmt.Println("Invalid operation")
		return
	}
	if operator != "*" && operator != "/" && operator != "+" && operator != "-" {
		fmt.Println("Invalid operation")
		return
	}

	ans := Count(input1, input2, operator)
	fmt.Println(ans)
}

func Count(a, b int, operator string) int {
	switch operator {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			os.Exit(1)
		}
		return a / b
	}
	return 0
}
