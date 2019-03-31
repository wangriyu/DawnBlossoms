package one

import (
	"math"
)

type UrlObj struct {
	url   string
	count int
}

type MinHeap struct {
	max  int
	len  int
	tree []UrlObj
}

func NewMinHeap(max int) *MinHeap {
	heap := &MinHeap{max: max, tree: make([]UrlObj, 1, max)}
	heap.tree[0] = UrlObj{count: math.MinInt64}
	return heap
}

func (heap *MinHeap) Push(x UrlObj) {
	if heap.len >= heap.max {
		heap.Pop()
	}
	heap.tree = append(heap.tree, x)
	heap.len++
	i := heap.len
	for ; heap.tree[i/2].count > x.count; i /= 2 {
		heap.tree[i] = heap.tree[i/2]
	}
	heap.tree[i] = x
}

func (heap *MinHeap) Top() int {
	if heap.len <= 0 {
		return heap.tree[0].count
	}
	return heap.tree[1].count
}

func (heap *MinHeap) Pop() UrlObj {
	if heap.len > 0 {
		min := heap.tree[1]
		last := heap.tree[heap.len]
		var i, child int
		for i = 1; i*2 <= heap.len; i = child {
			child = i * 2
			if child < heap.len && heap.tree[child+1].count < heap.tree[child].count {
				child++
			}
			if last.count > heap.tree[child].count {
				heap.tree[i] = heap.tree[child]
			} else {
				break
			}
		}
		heap.tree[i] = last
		heap.tree = heap.tree[:heap.len]
		heap.len--
		return min
	}
	return heap.tree[0]
}
