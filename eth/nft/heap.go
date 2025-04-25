package nft

import goethTypes "github.com/ethereum/go-ethereum/core/types"

type LogHeap []goethTypes.Log

func (h LogHeap) Len() int           { return len(h) }
func (h LogHeap) Less(i, j int) bool { return h[i].BlockNumber < h[j].BlockNumber }
func (h LogHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LogHeap) Push(x interface{}) {
	*h = append(*h, x.(goethTypes.Log))
}

func (h *LogHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
