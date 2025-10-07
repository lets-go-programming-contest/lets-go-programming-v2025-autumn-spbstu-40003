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

	var n, k int
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &k)

	for dep := 0; dep < n; dep++ {
		lo, hi := 15, 30
		for i := 0; i < k; i++ {
			var sign string
			var temp int
			fmt.Fscan(in, &sign, &temp)

			if sign == ">=" {
				if temp > lo {
					lo = temp
				}
			} else if sign == "<=" {
				if temp < hi {
					hi = temp
				}
			}
		}

		if lo <= hi {
			fmt.Fprintln(out, lo)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
