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
	_, errB := fmt.Scanln(&inputB)
	_, errOp := fmt.Scanln(&inputOp)

	if errA != nil {
		return 0, ErrInvalidFirstOperand
	}
	if errB != nil {
		return 0, ErrInvalidSecondOperand
	}
	if errOp != nil || len(inputOp) != 1 {
		return 0, ErrInvalidOperation
	}

	a, err1 := strconv.Atoi(inputA)
	b, err2 := strconv.Atoi(inputB)

	if err1 != nil {
		return 0, ErrInvalidFirstOperand
	}
	if err2 != nil {
		return 0, ErrInvalidSecondOperand
	}

	operator := inputOp[0]

	if operator == '/' && b == 0 {
		return 0, ErrDivisionByZero
	}

	switch operator {
	case '+':
		return float64(a + b), nil
	case '-':
		return float64(a - b), nil
	case '*':
		return float64(a * b), nil
	case '/':
		return float64(a) / float64(b), nil
	default:
		return 0, ErrInvalidOperation
	}
}
