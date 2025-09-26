package intheap

type IntHeap []int

func (h *IntHeap) Len() int { return len(*h) }

func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	if x, ok := x.(int); ok {
		*h = append(*h, x)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func (h *IntHeap) Peek() interface{} {
	return (*h)[0]
}