package intheap

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(first, second int) bool {
	return (*heap)[first] < (*heap)[second]
}

func (heap *IntHeap) Swap(first, second int) {
	(*heap)[first], (*heap)[second] = (*heap)[second], (*heap)[first]
}

func (heap *IntHeap) Push(object interface{}) {
	if val, ok := object.(int); ok {
		*heap = append(*heap, val)
	}
}

func (heap *IntHeap) Pop() interface{} {
	old := *heap
	oldLen := len(old)
	popped := old[oldLen-1]
	*heap = old[0 : oldLen-1]

	return popped
}
