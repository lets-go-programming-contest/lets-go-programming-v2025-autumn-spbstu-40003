package maxheap

type MaxIntHeap []int

func (heap *MaxIntHeap) Len() int {
	return len(*heap)
}

func (heap *MaxIntHeap) Less(i, j int) bool {
	return (*heap)[i] > (*heap)[j]
}

func (heap *MaxIntHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *MaxIntHeap) Push(data any) {
	value, ok := data.(int)
	if ok {
		*heap = append(*heap, value)
	}
}

func (heap *MaxIntHeap) Pop() any {
	old := *heap
	index := len(old)
	value := old[index-1]
	*heap = old[0 : index-1]

	return value
}
