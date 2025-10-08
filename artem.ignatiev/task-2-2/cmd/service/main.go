package main

import (
	"container/heap"
	"fmt"
)

const (
	errorValue = 0
	minDishes  = 1
	maxDishes  = 10000
	minRating  = -10000
	maxRating  = 10000
)

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	var n, k int

	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println(errorValue)
		return
	}

	if n < minDishes || n > maxDishes {
		fmt.Println(errorValue)
		return
	}

	dishes := make([]int, n)
	for i := range dishes {
		if _, err := fmt.Scan(&dishes[i]); err != nil {
			fmt.Println(errorValue)
			return
		}
		if dishes[i] < minRating || dishes[i] > maxRating {
			fmt.Println(errorValue)
			return
		}
	}

	if _, err := fmt.Scan(&k); err != nil {
		fmt.Println(errorValue)
		return
	}

	if k < minDishes || k > n {
		fmt.Println(errorValue)
		return
	}

	result, err := findKthPreference(dishes, k)
	if err != nil {
		fmt.Println(errorValue)
		return
	}

	fmt.Println(result)
}

func findKthPreference(dishes []int, k int) (int, error) {
	h := &MinHeap{}
	heap.Init(h)

	for _, rating := range dishes {
		heap.Push(h, rating)
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	if h.Len() == 0 {
		return errorValue, fmt.Errorf("empty heap")
	}

	return (*h)[0], nil
}
