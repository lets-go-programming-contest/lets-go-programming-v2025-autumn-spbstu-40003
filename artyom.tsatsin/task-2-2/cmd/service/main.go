package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrWrongCount  = errors.New("invalid number of menu items")
	ErrWrongRating = errors.New("invalid rating value")
	ErrWrongChoice = errors.New("invalid k value")
)

const (
	minRating = -10000
	maxRating = 10000
	minItems  = 1
)

type RatingHeap []int

func (h *RatingHeap) Len() int {
	return len(*h)
}

func (h *RatingHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *RatingHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *RatingHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *RatingHeap) Pop() interface{} {
	current := *h
	n := len(current)
	elem := current[n-1]
	*h = current[:n-1]

	return elem
}

func main() {
	var (
		dishCount int
		kSelect   int
	)

	_, err := fmt.Scan(&dishCount)
	if err != nil || dishCount < minItems || dishCount > maxRating {
		fmt.Println(ErrWrongCount)

		return
	}

	menuHeap, err := getRatings(dishCount, minRating, maxRating)
	if err != nil {
		fmt.Println(err)

		return
	}

	_, err = fmt.Scan(&kSelect)
	if err != nil || kSelect < 1 || kSelect > dishCount {
		fmt.Println(ErrWrongChoice)

		return
	}

	for kSelect > 0 {
		result, ok := heap.Pop(menuHeap).(int)
		if !ok {
			fmt.Println(-1)

			return
		}

		kSelect--
		if kSelect == 0 {
			fmt.Println(result)
		}
	}
}

func getRatings(count, minVal, maxVal int) (*RatingHeap, error) {
	data := &RatingHeap{}
	heap.Init(data)

	for range make([]struct{}, count) {
		var score int
		_, err := fmt.Scan(&score)

		if err != nil || score < minVal || score > maxVal {
			return nil, ErrWrongRating
		}

		heap.Push(data, score)
	}

	return data, nil
}
