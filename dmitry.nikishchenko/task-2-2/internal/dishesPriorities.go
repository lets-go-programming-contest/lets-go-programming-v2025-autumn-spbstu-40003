package dishesPriorities

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	errDishes      = errors.New("failed to read the number of dishes")
	errPriorities  = errors.New("failed to read the priorities")
	errPriorityNum = errors.New("failed to read priority num")
)

type DishesHeap []int

func (h DishesHeap) Len() int           { return len(h) }
func (h DishesHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h DishesHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *DishesHeap) Push(x any) {
	*h = append(*h, x.(int))
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
		n      int
		val    int
		k      int
		result int
	)

	if _, err := fmt.Scan(&n); err != nil {
		return errDishes
	}

	h := &DishesHeap{}
	heap.Init(h)

	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&val); err != nil {
			return errPriorities
		}
		heap.Push(h, val)
	}

	if _, err := fmt.Scan(&k); err != nil {
		return errPriorityNum
	}

	for i := 0; i < k; i++ {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)
	return nil
}
