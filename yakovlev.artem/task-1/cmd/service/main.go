package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	firstLine, _ := in.ReadString('\n')
	secondLine, _ := in.ReadString('\n')
	opLine, _ := in.ReadString('\n')

	firstLine = strings.TrimSpace(firstLine)
	secondLine = strings.TrimSpace(secondLine)
	opLine = strings.TrimSpace(opLine)

	a, err := strconv.Atoi(firstLine)
	if err != nil {
		fmt.Fprintln(out, "Invalid first operand")
		return
	}

	b, err := strconv.Atoi(secondLine)
	if err != nil {
		fmt.Fprintln(out, "Invalid second operand")
		return
	}

	switch opLine {
	case "+":
		fmt.Fprintln(out, a+b)
	case "-":
		fmt.Fprintln(out, a-b)
	case "*":
		fmt.Fprintln(out, a*b)
	case "/":
		if b == 0 {
			fmt.Fprintln(out, "Division by zero")
			return
		}
		fmt.Fprintln(out, a/b)
	default:
		fmt.Fprintln(out, "Invalid operation")
	}
}
