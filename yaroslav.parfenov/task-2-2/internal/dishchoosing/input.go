package dishchoosing

import (
	"container/heap"
	"fmt"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

const (
	errInvNumDishes      = "inv num dishes: %w"
	errInvSomeDish       = "inv some dish: %w"
	errInvOrdPerfectDish = "inv ord-perfect dish: %w"
)

func GetInput(dishStorage *ih.IntHeap) (int, error) {
	var numDishes int

	_, err := fmt.Scanln(&numDishes)
	if err != nil {
		return 0, fmt.Errorf(errInvNumDishes, err)
	}

	for range numDishes {
		var curDish int

		_, err = fmt.Scan(&curDish)
		if err != nil {
			return 0, fmt.Errorf(errInvSomeDish, err)
		}

		heap.Push(dishStorage, curDish)
	}

	var ordPerfectDish int

	_, err = fmt.Scanln(&ordPerfectDish)
	if err != nil {
		return 0, fmt.Errorf(errInvOrdPerfectDish, err)
	}

	return ordPerfectDish, nil
}
