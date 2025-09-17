package main

import "fmt"

func main() {
	var num1, num2, res int
	var operator rune

	_, error_num1 := fmt.Scanf("%d\n", &num1)

	if error_num1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, error_num2 := fmt.Scanf("%d\n", &num2)

	if error_num2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scanf("%c\n", &operator)

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
