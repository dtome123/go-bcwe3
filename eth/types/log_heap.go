package types

type LogHeap []Log

func (h LogHeap) Len() int           { return len(h) }
func (h LogHeap) Less(i, j int) bool { return h[i].BlockNumber < h[j].BlockNumber }
func (h LogHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LogHeap) Push(x any) {
	xLog, ok := x.(Log)

	if ok {
		*h = append(*h, xLog)
	}
}

func (h *LogHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
