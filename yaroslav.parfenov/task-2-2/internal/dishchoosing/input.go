package dishchoosing

import (
	"container/heap"
	"errors"
	"fmt"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

var (
	errInvNumDishes      = errors.New("inv num dishes")
	errInvSomeDish       = errors.New("inv some dish")
	errInvOrdPerfectDish = errors.New("inv ord-perfect dish")
)

func GetInput(dishStorage *ih.IntHeap) (int, error) {
	var numDishes int

	_, err := fmt.Scanln(&numDishes)
	if err != nil {
		return 0, errInvNumDishes
	}

	for range numDishes {
		var curDish int

		_, err = fmt.Scan(&curDish)
		if err != nil {
			return 0, errInvSomeDish
		}

		heap.Push(dishStorage, curDish)
	}

	var ordPerfectDish int

	_, err = fmt.Scanln(&ordPerfectDish)

	if err != nil {
		return 0, errInvOrdPerfectDish
	}

	return ordPerfectDish, nil
}
