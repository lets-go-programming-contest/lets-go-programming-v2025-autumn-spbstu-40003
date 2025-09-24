package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wedwincode/task-2-1/internal/temperature"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	if err := temperature.Solve(in, temperature.DefaultLowerBound, temperature.DefaultUpperBound); err != nil {
		fmt.Println(err)
	}
}
