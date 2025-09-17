package main

import "fmt"

func main() {
	var num1, num2 int
	var oper string

	_, err1 := fmt.Scan(&num1)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scan(&num2)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scan(&oper)

	switch oper[0] {
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
