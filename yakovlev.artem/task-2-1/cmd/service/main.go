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

	for { // читаем наборы данных до EOF
		var departments, employees int
		if _, err := fmt.Fscan(reader, &departments); err != nil {
			break // EOF — выходим
		}
		if _, err := fmt.Fscan(reader, &employees); err != nil {
			return
		}

		for range departments { // Go 1.22: 0..departments-1
			low, high := 15, 30

			for range employees { // 0..employees-1
				var op string
				var t int
				if _, err := fmt.Fscan(reader, &op, &t); err != nil {
					return
				}

				switch op {
				case ">=":
					if t > low {
						low = t
					}
				case "<=":
					if t < high {
						high = t
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
}
