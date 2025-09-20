package main

import (
	"fmt"
	"github.com/ArtttNik/task-1/internal"
)

func main() {
	result, err := internal.Calculate()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
