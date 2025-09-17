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

	_, errA := fmt.Scanln(&inputA)
	if errA != nil {
		return 0, ErrInvalidFirstOperand
	}
	_, errB := fmt.Scanln(&inputB)
	if errB != nil {
		return 0, ErrInvalidSecondOperand
	}
	_, errOp := fmt.Scanln(&inputOp)
	if errOp != nil || len(inputOp) != 1 {
		return 0, ErrInvalidOperation
	}

	a, err1 := strconv.Atoi(inputA)
	if err1 != nil {
		return 0, ErrInvalidFirstOperand
	}
	b, err2 := strconv.Atoi(inputB)
	if err2 != nil {
		return 0, ErrInvalidSecondOperand
	}

	operator := inputOp[0]
	switch operator {
	case '+':
		return float64(a + b), nil
	case '-':
		return float64(a - b), nil
	case '*':
		return float64(a * b), nil
	case '/':
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return float64(a) / float64(b), nil
	default:
		return 0, ErrInvalidOperation
	}
}
