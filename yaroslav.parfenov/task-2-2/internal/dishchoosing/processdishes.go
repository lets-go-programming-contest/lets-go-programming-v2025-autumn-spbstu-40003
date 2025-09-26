package dishchoosing

import (
	"container/heap"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

func ProcessDishes(dishStorage *ih.IntHeap, ordPerfectDish int) (int, error) {

	var curNumDish int
	for ordPerfectDish > 0 {
		curNumDish = heap.Pop(dishStorage).(int)
		ordPerfectDish--
	}

	return curNumDish, nil
}
