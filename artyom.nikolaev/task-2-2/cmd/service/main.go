package main

import (
	"fmt"

	"github.com/ArtttNik/task-2-2/internal"
)

func main() {
	err := internal.FindKDish()
	if err != nil {
		fmt.Println(err)

		return
	}
}
