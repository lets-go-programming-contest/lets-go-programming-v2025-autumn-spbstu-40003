package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func readLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if errors.Is(err, io.EOF) {
		if len(s) == 0 {
			return "", err
		}
		err = nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimRight(s, "\r\n"), nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer func() {
		if err := out.Flush(); err != nil {
			// не печатаем в stdout, чтобы не ломать протокол вывода
			fmt.Fprintln(os.Stderr, "flush error:", err)
		}
	}()

	firstLine, err := readLine(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read first line:", err)
		return
	}
	secondLine, err := readLine(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read second line:", err)
		return
	}
	opLine, err := readLine(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read operation:", err)
		return
	}

	a, err := strconv.Atoi(strings.TrimSpace(firstLine))
	if err != nil {
		fmt.Fprintln(out, "Invalid first operand")
		return
	}
	b, err := strconv.Atoi(strings.TrimSpace(secondLine))
	if err != nil {
		fmt.Fprintln(out, "Invalid second operand")
		return
	}

	switch strings.TrimSpace(opLine) {
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
