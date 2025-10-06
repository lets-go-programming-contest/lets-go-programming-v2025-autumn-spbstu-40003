package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int // количество отделов
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	for dep := 0; dep < n; dep++ {
		var k int // сотрудников в отделе
		fmt.Fscan(in, &k)

		lo, hi := 15, 30
		for i := 0; i < k; i++ {
			var op string
			var v int
			fmt.Fscan(in, &op, &v)

			switch op {
			case ">=":
				if v > lo {
					lo = v
				}
			case "<=":
				if v < hi {
					hi = v
				}
			default:
				// игнорируем неожиданные токены
			}

			if lo <= hi {
				fmt.Fprintln(out, lo)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}
