package main

import (
	"fmt"

	dc "github.com/gituser549/task-2-2/internal/dishchoosing"
	ih "github.com/gituser549/task-2-2/internal/intheap"
)

func main() {
	var dishStorage ih.IntHeap

	ordPerfectDish, err := dc.GetInput(&dishStorage)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(dc.ProcessDishes(&dishStorage, ordPerfectDish))
}
