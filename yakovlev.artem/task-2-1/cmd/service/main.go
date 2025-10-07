package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer func() { _ = writer.Flush() }()

	var departments, employees int
	if _, err := fmt.Fscan(reader, &departments); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &employees); err != nil {
		return
	}

	for range departments { // Go 1.22: интерация по 0..departments-1
		low, high := 15, 30

		for range employees { // 0..employees-1
			var sign string
			var temp int

			if _, err := fmt.Fscan(reader, &sign, &temp); err != nil {
				return
			}

			switch sign {
			case ">=":
				if temp > low {
					low = temp
				}
			case "<=":
				if temp < high {
					high = temp
				}
			}

			if low <= high {
				_, _ = fmt.Fprintln(writer, low)
			} else {
				_, _ = fmt.Fprintln(writer, -1)
			}
		}
	}
}
