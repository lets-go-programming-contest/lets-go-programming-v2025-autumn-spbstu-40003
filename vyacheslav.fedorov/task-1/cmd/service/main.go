package main

import "fmt"

func main() {
	var (
		firNum, secNum int
		operation      string
	)

	_, err := fmt.Scan(&firNum)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&secNum)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(firNum + secNum)
	case "-":
		fmt.Println(firNum - secNum)
	case "*":
		fmt.Println(firNum * secNum)
	case "/":
		if secNum == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firNum / secNum)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
