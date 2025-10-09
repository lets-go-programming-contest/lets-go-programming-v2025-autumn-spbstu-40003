package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrInvalidDish       = errors.New("incorrect input count of dish")
	ErrInvalidValue      = errors.New("incorrect dish value")
	ErrInvalidPreference = errors.New("incorrect preference of dish value")
)

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(i, j int) bool {
	return (*heap)[i] > (*heap)[j]
}

func (heap *IntHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *IntHeap) Push(value interface{}) {
	if intValue, ok := value.(int); ok {
		*heap = append(*heap, intValue)
	}
}

func (heap *IntHeap) Pop() interface{} {
	currentHeap := *heap
	length := len(currentHeap)
	lastElement := currentHeap[length-1]
	*heap = currentHeap[0 : length-1]

	return lastElement
}

func main() {
	var (
		countOfDish      int
		preferenceOfDish int
	)

	const (
		MinValue  = -10000
		MinDishes = 1
		MaxValue  = 10000
	)

	_, err := fmt.Scan(&countOfDish)
	if err != nil || countOfDish < MinDishes || countOfDish > MaxValue {
		fmt.Println(ErrInvalidDish)

		return
	}

	heapData, err := readAndValidateDishes(countOfDish, MinValue, MaxValue)
	if err != nil {
		fmt.Println(err)

		return
	}

	_, err = fmt.Scan(&preferenceOfDish)
	if err != nil || preferenceOfDish < 1 || preferenceOfDish > countOfDish {
		fmt.Println(ErrInvalidPreference)

		return
	}

	for preferenceOfDish > 0 {
		result, isInt := heap.Pop(heapData).(int)
		if !isInt {
			fmt.Println(-1)

			return
		}

		preferenceOfDish--

		if preferenceOfDish == 0 {
			fmt.Println(result)
		}
	}
}

func readAndValidateDishes(count, minValue, maxValue int) (*IntHeap, error) {
	heapData := &IntHeap{}
	heap.Init(heapData)

	for range count {
		var value int
		_, err := fmt.Scan(&value)

		if err != nil || value < minValue || value > maxValue {
			return nil, ErrInvalidValue
		}

		heap.Push(heapData, value)
	}

	return heapData, nil
}
