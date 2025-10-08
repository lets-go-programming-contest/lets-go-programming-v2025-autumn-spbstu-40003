package main

import (
	"container/heap"
	"errors"
	"fmt"
)

const errorValue = 0

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

func main() {
	var (
		dishCount      int
		dishPreference int
	)

	if _, err := fmt.Scan(&dishCount); err != nil {
		handleError(fmt.Errorf("%w: %v", ErrReadInput, err))

		return
	}

	if dishCount < 1 || dishCount > 10000 {
		handleError(ErrDishCountRange)

		return
	}

	dishRatings := make([]int, dishCount)
	for i := range dishRatings {
		if _, err := fmt.Scan(&dishRatings[i]); err != nil {
			handleError(fmt.Errorf("%w: %v", ErrReadInput, err))

			return
		}

		if dishRatings[i] < -10000 || dishRatings[i] > 10000 {
			handleError(ErrRatingRange)

			return
		}
	}

	if _, err := fmt.Scan(&dishPreference); err != nil {
		handleError(fmt.Errorf("%w: %v", ErrReadInput, err))

		return
	}

	if dishPreference < 1 || dishPreference > dishCount {
		handleError(ErrPreferenceRange)

		return
	}

	result := getPreference(dishRatings, dishPreference)
	fmt.Println(result)
}

func getPreference(dishRatings []int, preference int) int {
	ratingHeap := &IntHeap{}
	heap.Init(ratingHeap)

	for _, rating := range dishRatings {
		heap.Push(ratingHeap, rating)

		if ratingHeap.Len() > preference {
			heap.Pop(ratingHeap)
		}
	}

	return (*ratingHeap)[0]
}

func handleError(err error) {
	_ = err
	fmt.Println(errorValue)
}

var (
	ErrReadInput       = errors.New("failed to read input")
	ErrDishCountRange  = errors.New("dish count out of range")
	ErrRatingRange     = errors.New("dish rating out of range")
	ErrPreferenceRange = errors.New("preference out of range")
)
