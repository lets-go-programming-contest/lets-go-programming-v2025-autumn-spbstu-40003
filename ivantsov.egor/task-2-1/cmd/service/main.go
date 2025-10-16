package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const minTemp = 15
	const maxTemp = 30

	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)

	// Читаем количество отделов N и сотрудников K
	var n, k int
	in.Scan()
	fmt.Sscan(in.Text(), &n)
	in.Scan()
	fmt.Sscan(in.Text(), &k)

	for i := 0; i < n; i++ {
		minBound := minTemp
		maxBound := maxTemp

		for j := 0; j < k; j++ {
			in.Scan()
			sign := in.Text()
			in.Scan()
			var t int
			fmt.Sscan(in.Text(), &t)

			if sign == ">=" {
				if t > minBound {
					minBound = t
				}
			} else if sign == "<=" {
				if t < maxBound {
					maxBound = t
				}
			}
		}

		if minBound <= maxBound {
			fmt.Println(maxBound)
		} else {
			fmt.Println(-1)
		}
	}
}
