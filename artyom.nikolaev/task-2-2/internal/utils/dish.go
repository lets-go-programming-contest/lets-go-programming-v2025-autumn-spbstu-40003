package utils

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrInputFail     = errors.New("input error")
	ErrInvalidAmount = errors.New("invalid amount")
)

const (
	MaxDishes = 10000
	MinValue  = -10000
	MaxValue  = 10000
)

type intHeap []int

func (h *intHeap) Len() int {
	return len(*h)
}

func (h *intHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *intHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *intHeap) Push(x any) {
	v, isInt := x.(int)
	if isInt {
		*h = append(*h, v)
	}
}

func (h *intHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func readAndParseInput() ([]int, int, error) {
	var dishesCount int

	_, err := fmt.Scan(&dishesCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading dishes count: %w", ErrInputFail)
	}

	if dishesCount <= 0 || dishesCount > MaxDishes {
		return nil, 0, fmt.Errorf("error of the range of dishes count: %w", ErrInvalidAmount)
	}

	dishesSlice := make([]int, dishesCount)

	for index := range dishesSlice {
		var cost int

		_, err := fmt.Scan(&cost)
		if err != nil {
			return nil, 0, fmt.Errorf("error reading dish cost: %w", ErrInputFail)
		}

		if cost < MinValue || cost > MaxValue {
			return nil, 0, fmt.Errorf("error of the range of dish cost: %w", ErrInvalidAmount)
		}

		dishesSlice[index] = cost
	}

	var need int

	_, err = fmt.Scan(&need)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading k value: %w", ErrInputFail)
	}

	if need <= 0 || need > dishesCount {
		return nil, 0, fmt.Errorf("error of the range of k value: %w", ErrInvalidAmount)
	}

	return dishesSlice, need, nil
}

func findKthLargest(prices []int, k int) int {
	minHeap := make(intHeap, k)
	copy(minHeap, prices[:k])
	heap.Init(&minHeap)

	for _, price := range prices[k:] {
		if price > minHeap[0] {
			minHeap[0] = price
			heap.Fix(&minHeap, 0)
		}
	}

	return minHeap[0]
}

func FindKDish() error {
	dishes, k, err := readAndParseInput()
	if err != nil {
		return fmt.Errorf("error in FindKDish: %w", err)
	}

	result := findKthLargest(dishes, k)
	fmt.Println(result)

	return nil
}
