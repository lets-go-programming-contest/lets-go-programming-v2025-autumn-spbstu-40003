package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	firstInput, _ := reader.ReadString('\n')
	firstInput = strings.TrimSpace(firstInput)
	firstOperand, err := strconv.Atoi(firstInput)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	secondInput, _ := reader.ReadString('\n')
	secondInput = strings.TrimSpace(secondInput)
	secondOperand, err := strconv.Atoi(secondInput)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	opInput, _ := reader.ReadString('\n')
	opInput = strings.TrimSpace(opInput)

	switch opInput {
	case "+":
		fmt.Println(firstOperand + secondOperand)
	case "-":
		fmt.Println(firstOperand - secondOperand)
	case "*":
		fmt.Println(firstOperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(firstOperand / secondOperand)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
