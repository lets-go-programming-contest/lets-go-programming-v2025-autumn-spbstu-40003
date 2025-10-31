package ratingheap

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
