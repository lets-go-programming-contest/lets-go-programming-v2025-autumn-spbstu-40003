package main

import (
	"fmt"
	"strconv"
)

func inputOperand(opd *int, opdNumber string) bool {
	var input string
	var err error
	fmt.Scanln(&input)
	*opd, err = strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid %s operand\n", opdNumber)
		return false
	}
	return true
}

func main() {
	var opd1, opd2 int
	flag := inputOperand(&opd1, "first")
	if !flag {
		return
	}
	flag = inputOperand(&opd2, "second")
	if !flag {
		return
	}

	var opn string
	fmt.Scanln(&opn)

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
		result = int(opd1 / opd2)
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)
}