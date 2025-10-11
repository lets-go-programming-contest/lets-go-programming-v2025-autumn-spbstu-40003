package main

import "fmt"

func plus (a float64, b float64) float64 {
	return a + b
}

func minus (a float64, b float64) float64 {
	return a - b
}

func multiply (a float64, b float64) float64 {
	return a * b
}

func divide (a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("Division by zero")
	}
	return a / b, nil
}

func readCommands () (float64, float64, string, error) {
	var a, b float64
	var operation string

	_, err1 := fmt.Scanln(&a)
	_, err2 := fmt.Scanln(&b)
	_, err3 := fmt.Scanln(&operation)

	if err1 != nil {
		return -1, -1, "", fmt.Errorf("Invalid first operand")
	}

	if err2 != nil {
		return -1, -1, "", fmt.Errorf("Invalid second operand")
	}

	if err3 != nil {
		return -1, -1, "", fmt.Errorf("Invalid operation")
	}

	return a, b, operation, nil
}

func executeCommand (a float64, b float64, operation string) (float64, error) {
	switch operation {
		case "+":
			return plus (a, b), nil
		case "-":
			return minus (a, b), nil
		case "*":
			return multiply (a, b), nil
		case "/":
			return divide (a, b)
		default:
			return -1, fmt.Errorf("Invalid operation")
	}
}

func main() {
	var a, b float64
	var operation string
	var result float64

	a, b, operation, err := readCommands()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, errr := executeCommand (a, b, operation)

	if errr != nil {
		fmt.Println (errr.Error())
		return
	}

	if result == float64(int(result)) {
		fmt.Println(int(result))
	} else {
		fmt.Println(result)
	}
}
