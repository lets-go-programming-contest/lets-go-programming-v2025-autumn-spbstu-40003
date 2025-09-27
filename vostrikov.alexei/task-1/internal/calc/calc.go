package calc

import (
	"errors"
	"fmt"
	"strconv"
)

func input() (int, int, string, error) {
	var (
		str1   string
		str2   string
		symbol string

		num1, num2 int
	)

	_, err := fmt.Scanln(&str1)
	if err != nil {
		return 0, 0, "", errors.New("Invalid first operand")
	}
	num1, err = strconv.Atoi(str1)
	if err != nil {
		return 0, 0, "", errors.New("Invalid first operand")
	}

	_, err = fmt.Scanln(&str2)
	if err != nil {
		return 0, 0, "", errors.New("Invalid second operand")
	}
	num2, err = strconv.Atoi(str2)
	if err != nil {
		return 0, 0, "", errors.New("Invalid second operand")
	}

	_, err = fmt.Scanln(&symbol)
	if err != nil {
		return 0, 0, "", errors.New("Problem in symbol")
	}
	if !validate(symbol) {
		return 0, 0, "", errors.New("Invalid operation")
	}
	if symbol == "/" && str2 == "0" {
		return 0, 0, "", errors.New("Division by zero")
	}

	return num1, num2, symbol, nil
}

func Calculate() (float32, error) {
	num1, num2, symbol, errCode := input()

	if errCode != nil {
		return 0, errCode
	}

	var res float32

	switch symbol {
	case "+":
		res = float32(num1 + num2)
	case "-":
		res = float32(num1 - num2)
	case "*":
		res = float32(num1 * num2)
	case "/":
		res = float32(num1) / float32(num2)
	}
	return res, nil
}

func validate(str string) bool {
	var correct = [4]string{"-", "+", "*", "/"}

	for i := 0; i < len(correct); i++ {
		if str == correct[i] {
			return true
		}
	}
	return false
}
