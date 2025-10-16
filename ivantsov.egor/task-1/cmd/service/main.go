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
	firstInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	firstInput = strings.TrimSpace(firstInput)
	firstOperand, err := strconv.Atoi(firstInput)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	secondInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	secondInput = strings.TrimSpace(secondInput)
	secondOperand, err := strconv.Atoi(secondInput)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	opInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}
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
			return
		}
		fmt.Println(firstOperand / secondOperand)
	default:
		fmt.Println("Invalid operation")
	}
}
