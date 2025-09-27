package main

import "fmt"

func main() {

	var (
		firstNumber, secondNumber int
		operation                 string
	)

	_, err := fmt.Scan(&firstNumber)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondNumber)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	result, errMsg := calculateNumbers(firstNumber, secondNumber, operation)

	if errMsg != "" {
		fmt.Println(errMsg)
		return
	}

	fmt.Println(result)
}

func calculateNumbers(firstNumber, secondNumber int, operation string) (int, string) {
	switch operation {
	case "+":
		return firstNumber + secondNumber, ""
	case "-":
		return firstNumber - secondNumber, ""
	case "/":
		if secondNumber == 0 {
			return 0, "Division by zero"
		} else {
			return firstNumber / secondNumber, ""
		}
	case "*":
		return firstNumber * secondNumber, ""
	default:
		return 0, "Invalid operation"
	}
}
