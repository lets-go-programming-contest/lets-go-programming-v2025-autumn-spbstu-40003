package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer func() {
		_ = out.Flush()
	}()

	var numDepartments, numEmployees int
	if _, err := fmt.Fscan(in, &numDepartments); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &numEmployees); err != nil {
		return
	}

	for depIndex := 0; depIndex < numDepartments; depIndex++ {
		low, high := 15, 30

		for empIndex := 0; empIndex < numEmployees; empIndex++ {
			var sign string
			var temp int
			if _, err := fmt.Fscan(in, &sign, &temp); err != nil {
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
		}

		if low <= high {
			_, _ = fmt.Fprintln(out, low)
		} else {
			_, _ = fmt.Fprintln(out, -1)
		}
	}
}
