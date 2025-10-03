package main

import (
	"errors"
	"fmt"
)

func main() {
	var (
		first, second       int
		operation           string
		ErrDivByZero        = errors.New("Division by zero")
		ErrInvalidfirst     = errors.New("Invalid first operand")
		ErrInvalidsecond    = errors.New("Invalid second operand")
		ErrInvalidoperation = errors.New("Invalid operation")
	)

	_, err := fmt.Scan(&first)
	if err != nil {
		fmt.Println(ErrInvalidfirst)
		return
	}
	_, err = fmt.Scan(&second)
	if err != nil {
		fmt.Println(ErrInvalidsecond)
		return
	}
	_, err = fmt.Scan(&operation)

	switch operation {
	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	case "/":
		if second == 0 {
			fmt.Println(ErrDivByZero)
			return
		} else {
			fmt.Println(first / second)
		}
	default:
		fmt.Println(ErrInvalidoperation)

	}
}
