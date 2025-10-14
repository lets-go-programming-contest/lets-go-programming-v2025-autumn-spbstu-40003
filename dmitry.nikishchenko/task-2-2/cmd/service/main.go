package main

import (
	"fmt"

	dishesPriorities "github.com/d1mene/task-2-2/internal"
)

func main() {
	err := dishesPriorities.PickBestDish()
	if err != nil {
		fmt.Println(err)
	}
}
