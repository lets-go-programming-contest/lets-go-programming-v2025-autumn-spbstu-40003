package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func rl(r *bufio.Reader) string {
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

func main() {
	in := bufio.NewReader(os.Stdin)

	a, err := strconv.Atoi(rl(in))
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	b, err := strconv.Atoi(rl(in))
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	op := rl(in)
	if len(op) != 1 || !strings.Contains("+-*/", op) {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(a / b)
	}
}
