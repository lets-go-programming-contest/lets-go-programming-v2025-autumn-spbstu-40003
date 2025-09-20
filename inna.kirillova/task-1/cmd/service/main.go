package main

import "fmt"

func main() {
	var a, b int
	var operation string

	_, errorOfInvalidFirstOperand := fmt.Scanln(&a)
	if errorOfInvalidFirstOperand != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, errorOfInvalidSecondOperand := fmt.Scanln(&b)
	if errorOfInvalidSecondOperand != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, errorOfInvalidOperation := fmt.Scanln(&operation)
	if errorOfInvalidOperation != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
