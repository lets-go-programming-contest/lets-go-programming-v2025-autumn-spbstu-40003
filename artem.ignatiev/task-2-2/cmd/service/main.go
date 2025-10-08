package main

import (
	"container/heap"
	"errors"
	"fmt"
)

const (
	errorValue = 0
	minDishes  = 1
	maxDishes  = 10000
	minRating  = -10000
	maxRating  = 10000
)

var errEmptyHeap = errors.New("empty heap")

type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	var numDishes, preferenceIndex int

	if _, err := fmt.Scan(&numDishes); err != nil {
		fmt.Println(errorValue)

		return
	}

	if numDishes < minDishes || numDishes > maxDishes {
		fmt.Println(errorValue)

		return
	}

	dishRatings := make([]int, numDishes)
	for index := range dishRatings {
		if _, err := fmt.Scan(&dishRatings[index]); err != nil {
			fmt.Println(errorValue)

			return
		}

		if dishRatings[index] < minRating || dishRatings[index] > maxRating {
			fmt.Println(errorValue)

			return
		}
	}

	if _, err := fmt.Scan(&preferenceIndex); err != nil {
		fmt.Println(errorValue)

		return
	}

	if preferenceIndex < minDishes || preferenceIndex > numDishes {
		fmt.Println(errorValue)

		return
	}

	result, err := findKthPreference(dishRatings, preferenceIndex)
	if err != nil {
		fmt.Println(errorValue)

		return
	}

	fmt.Println(result)
}

func findKthPreference(dishRatings []int, preferenceIndex int) (int, error) {
	heapInstance := &MinHeap{}
	heap.Init(heapInstance)

	for _, rating := range dishRatings {
		heap.Push(heapInstance, rating)

		if heapInstance.Len() > preferenceIndex {
			heap.Pop(heapInstance)
		}
	}

	if heapInstance.Len() == 0 {
		return errorValue, errEmptyHeap
	}

	return (*heapInstance)[0], nil
}
