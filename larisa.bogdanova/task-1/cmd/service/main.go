package main

import "fmt"

func main() {
	var (
		num1, num2 int
		oper       rune
	)
	_, err1 := fmt.Scanln(&num1)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scanln(&num2)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err3 := fmt.Scanf("%c", &oper)
	if err3 != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch oper {
	case '+':
		fmt.Println(num1 + num2)
		return
	case '-':
		fmt.Println(num1 - num2)
		return
	case '*':
		fmt.Println(num1 * num2)
		return
	case '/':
		if num2 != 0 {
			fmt.Println(num1 / num2)
		} else {
			fmt.Println("Division by zero")
			return
		}
	default:
		fmt.Println("Invalid operation")
	}
}
