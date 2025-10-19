package dishespriorities

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	errDishes      = errors.New("failed to read the number of dishes")
	errPriorities  = errors.New("failed to read the priorities")
	errPriorityNum = errors.New("failed to read priority num")
	errFormat      = errors.New("failed to convert heap element to int")
	errHeapRange   = errors.New("priority num out of range")
)

type DishesHeap []int

func (h *DishesHeap) Len() int           { return len(*h) }
func (h *DishesHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *DishesHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *DishesHeap) Push(x any) {
	val, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, val)
}

func (h *DishesHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func PickBestDish() error {
	var (
		numberOfDishes int
		val            int
		priorityNum    int
		result         int
	)

	if _, err := fmt.Scan(&numberOfDishes); err != nil {
		return errors.Join(errDishes, err)
	}

	priorityHeap := &DishesHeap{}
	heap.Init(priorityHeap)

	for range numberOfDishes {
		if _, err := fmt.Scan(&val); err != nil {
			return errors.Join(errPriorities, err)
		}

		heap.Push(priorityHeap, val)
	}

	if _, err := fmt.Scan(&priorityNum); err != nil {
		return errors.Join(errPriorityNum, err)
	}

	if priorityNum > priorityHeap.Len() || priorityNum < 1 {
		return errHeapRange
	}

	for range priorityNum {
		val := heap.Pop(priorityHeap)
		if valInt, ok := val.(int); ok {
			result = valInt
		} else {
			return errFormat
		}
	}

	fmt.Println(result)

	return nil
}
