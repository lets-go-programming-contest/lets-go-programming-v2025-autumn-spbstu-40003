package main

import (
	"fmt"
	"strconv"
)

func inputOperand(opd *int, opdNumber string) error {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("Invalid %s operand", opdNumber)
	}

	*opd, err = strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("Invalid %s operand", opdNumber)
	}
	return nil
}

func main() {
	var opd1, opd2 int

	err := inputOperand(&opd1, "first");
	if err != nil {
		fmt.Println(err)
		return
	}

	err = inputOperand(&opd2, "second");
	if err != nil {
		fmt.Println(err)
		return
	}

	var opn string
	_, err = fmt.Scanln(&opn)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	var result int
	switch opn {
	case "+":
		result = opd1 + opd2
	case "-":
		result = opd1 - opd2
	case "*":
		result = opd1 * opd2
	case "/":
		if opd2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = opd1 / opd2
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)
}
