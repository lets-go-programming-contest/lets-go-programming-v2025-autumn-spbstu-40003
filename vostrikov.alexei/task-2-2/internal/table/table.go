package table

import (
	"container/heap"
	"errors"
	"fmt"

	myHeap "github.com/lexachanskii/task-2-2/internal/heap"
)

var (
	errOutOfRange    = errors.New("favourite dish number is bigger than dishes list length")
	errDifferentTypr = errors.New("unexpected type")
)

func inputHeap(myheap *myHeap.HeapInt) (int, error) {
	var (
		dishes    int
		favourite int
		dishNum   int
	)

	if _, err := fmt.Scan(&dishes); err != nil {
		return 0, fmt.Errorf("error while reading dishes count: %w", err)
	}

	for range dishes {
		if _, err := fmt.Scan(&dishNum); err != nil {
			return 0, fmt.Errorf("error while reading dishes list: %w", err)
		}

		heap.Push(myheap, dishNum)
	}

	if _, err := fmt.Scan(&favourite); err != nil {
		return 0, fmt.Errorf("error while reading favourite dish: %w", err)
	}

	return favourite, nil
}

func getFavourite(favourite int, myheap *myHeap.HeapInt) (int, error) {
	if myheap.Len() < favourite {
		return -1, fmt.Errorf("lenth error %w", errOutOfRange)
	}

	for range favourite - 1 {
		heap.Pop(myheap)
	}

	num, ok := heap.Pop(myheap).(int)
	if !ok {
		return -1, fmt.Errorf("error while popping value %w", errDifferentTypr)
	}

	return num, nil
}

func Table() (int, error) {
	var heap myHeap.HeapInt

	favourite, err := inputHeap(&heap)
	if err != nil {
		return -1, fmt.Errorf("error in heap input %w", err)
	}

	bestDish, err := getFavourite(favourite, &heap)
	if err != nil {
		return -1, fmt.Errorf("error while getting favourite dish %w", err)
	}

	return bestDish, nil
}
