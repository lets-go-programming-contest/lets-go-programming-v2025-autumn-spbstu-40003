package dishchoosing

import (
	"container/heap"

	ih "github.com/gituser549/task-2-2/internal/intheap"
)

func ProcessDishes() (int, error) {
	var dishStorage ih.IntHeap
	ordPerfectDish, err := getInput(&dishStorage)

	if err != nil {
		return 0, err
	}

	var curNumDish int
	for ordPerfectDish > 0 {
		curNumDish = heap.Pop(&dishStorage).(int)
		ordPerfectDish--
	}

	return curNumDish, nil
}
