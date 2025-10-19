package ownheap

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	data, ok := x.(int)
	if ok {
		*h = append(*h, data)
	}
}

func (h *IntHeap) Pop() any {
	old := *h
	length := len(old)
	data := old[length-1]
	*h = old[0 : length-1]

	return data
}
