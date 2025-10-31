package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/Artem-Hack/task-2-2/internal/ratingheap"
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

func main() {
	var (
		dishCount int
		kSelect   int
	)

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Printf("%v: %v\n", ErrWrongCount, err)

		return
	}

	if dishCount < minItems || dishCount > maxRating {
		fmt.Println(ErrWrongCount)

		return
	}

	menuHeap, err := getRatings(dishCount, minRating, maxRating)
	if err != nil {
		fmt.Println(err)

		return
	}

	_, err = fmt.Scan(&kSelect)
	if err != nil {
		fmt.Printf("%v: %v\n", ErrWrongChoice, err)

		return
	}

	if kSelect < 1 || kSelect > dishCount {
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

func getRatings(count, minVal, maxVal int) (*ratingheap.RatingHeap, error) {
	data := &ratingheap.RatingHeap{}
	heap.Init(data)

	for range count {
		var score int
		_, err := fmt.Scan(&score)

		if err != nil {
			return nil, fmt.Errorf("%w", ErrWrongRating)
		}

		if score < minVal || score > maxVal {
			return nil, ErrWrongRating
		}

		heap.Push(data, score)
	}

	return data, nil
}
