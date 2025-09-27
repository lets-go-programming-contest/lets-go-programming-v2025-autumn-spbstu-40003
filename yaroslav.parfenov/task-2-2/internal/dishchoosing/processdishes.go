package dishchoosing

import (
	"container/heap"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

func ProcessDishes(dishStorage *ih.IntHeap, ordPerfectDish int) int {
	var ansNumDish int

	for range ordPerfectDish - 1 {
		heap.Pop(dishStorage)
	}

	ansNumDish, ok := heap.Pop(dishStorage).(int)

	if !ok {
		return -1
	}

	return ansNumDish
}
