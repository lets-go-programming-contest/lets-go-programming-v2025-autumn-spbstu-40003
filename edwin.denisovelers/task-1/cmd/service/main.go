package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	errDivByZero        = errors.New("Division by zero")
	errBadOp            = errors.New("Invalid operation")
	errBadFirstOperand  = errors.New("Invalid first operand")
	errBadSecondOperand = errors.New("Invalid second operand")
)

func calculate(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errDivByZero
		}
		return a / b, nil
	default:
		return 0, errBadOp
	}
}

func readValue[T int | string](r io.Reader, val *T, errorToShow error) error {
	_, err := fmt.Fscanln(r, val)
	if err != nil {
		return errorToShow
	}
	return nil
}

func main() {
	in := bufio.NewReader(os.Stdin)

	var (
		a, b int
		op   string
	)

	if err := readValue(in, &a, errBadFirstOperand); err != nil {
		fmt.Println(err)
		return
	}
	if err := readValue(in, &b, errBadSecondOperand); err != nil {
		fmt.Println(err)
		return
	}
	if err := readValue(in, &op, errBadOp); err != nil {
		fmt.Println(err)
		return
	}

	res, err := calculate(a, b, op)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
