package main

import "fmt"

func main() {
	var (
		num1, num2, res int
		operator        rune
	)

	_, err := fmt.Scanln(&num1)

	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&num2)

	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scanf("%c", &operator)

	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case '+':
		res = num1 + num2
	case '-':
		res = num1 - num2
	case '*':
		res = num1 * num2
	case '/':
		if num2 != 0 {
			res = num1 / num2
		} else {
			fmt.Println("Division by zero")
			return
		}
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(res)
}
