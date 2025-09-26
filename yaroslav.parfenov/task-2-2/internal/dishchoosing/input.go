package dishchoosing

import (
	"container/heap"
	"fmt"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

func getInput(dishStorage *ih.IntHeap) (int, error) {
	var numDishes int

	_, err := fmt.Scanln(&numDishes)

	if err != nil {
		return 0, err
	}

	for numDishes > 0 {
		var curDish int

		_, err = fmt.Scan(&curDish)
		if err != nil {
			return 0, err
		}

		heap.Push(dishStorage, curDish)

		numDishes--
	}

	var ordPerfectDish int

	_, err = fmt.Scanln(&ordPerfectDish)

	if err != nil {
		return 0, err
	}

	return ordPerfectDish, nil
}
