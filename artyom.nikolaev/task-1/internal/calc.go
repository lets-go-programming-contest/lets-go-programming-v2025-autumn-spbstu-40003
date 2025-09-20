package internal

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidFirstOperand  = errors.New("Invalid first operand")
	ErrInvalidSecondOperand = errors.New("Invalid second operand")
	ErrInvalidOperation     = errors.New("Invalid operation")
	ErrDivisionByZero       = errors.New("Division by zero")
)

func Calculate() (float64, error) {
	var inputA, inputB, inputOp string

	if _, err := fmt.Scanln(&inputA); err != nil {
		return 0, ErrInvalidFirstOperand
	}
	a, err := strconv.Atoi(inputA)
	if err != nil {
		return 0, ErrInvalidFirstOperand
	}

	if _, err := fmt.Scanln(&inputB); err != nil {
		return 0, ErrInvalidSecondOperand
	}
	b, err := strconv.Atoi(inputB)
	if err != nil {
		return 0, ErrInvalidSecondOperand
	}

	if _, err := fmt.Scanln(&inputOp); err != nil {
		return 0, ErrInvalidOperation
	}

	switch inputOp {
	case "+":
		return float64(a + b), nil
	case "-":
		return float64(a - b), nil
	case "*":
		return float64(a * b), nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return float64(a) / float64(b), nil
	default:
		return 0, ErrInvalidOperation
	}
}
