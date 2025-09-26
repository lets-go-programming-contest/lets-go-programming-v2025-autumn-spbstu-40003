package main

import (
	"fmt"

	"github.com/ArtttNik/task-2-1/internal"
)

func main() {
	err := internal.Temp()
	if err != nil {
		fmt.Println(err)

		return
	}
}
