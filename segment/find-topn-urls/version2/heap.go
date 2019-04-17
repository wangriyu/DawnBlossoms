package version2

import (
	"math"
)

type UrlObj struct {
	Url   string
	Count int
}

type MinHeap struct {
	max  int
	len  int
	tree []UrlObj
}

func NewMinHeap(max int) *MinHeap {
	heap := &MinHeap{max: max, tree: make([]UrlObj, 1, max+1)}
	heap.tree[0] = UrlObj{Count: math.MinInt64}
	return heap
}

func (heap *MinHeap) Push(x UrlObj) {
	if heap.len >= heap.max {
		var i, child int
		for i = 1; i*2 <= heap.len; i = child {
			child = i * 2
			if child < heap.len && heap.tree[child+1].Count < heap.tree[child].Count {
				child++
			}
			if x.Count > heap.tree[child].Count {
				heap.tree[i] = heap.tree[child]
			} else {
				break
			}
		}
		heap.tree[i] = x
	} else {
		heap.tree = append(heap.tree, x)
		heap.len++
		i := heap.len
		for ; heap.tree[i/2].Count > x.Count; i /= 2 {
			heap.tree[i] = heap.tree[i/2]
		}
		heap.tree[i] = x
	}
}

func (heap *MinHeap) Top() int {
	if heap.len <= 0 {
		return heap.tree[0].Count
	}
	return heap.tree[1].Count
}

func (heap *MinHeap) PopToList() UrlObj {
	if heap.len > 0 {
		min := heap.tree[1]
		last := heap.tree[heap.len]
		var i, child int
		for i = 1; i*2 <= heap.len; i = child {
			child = i * 2
			if child < heap.len && heap.tree[child+1].Count < heap.tree[child].Count {
				child++
			}
			if last.Count > heap.tree[child].Count {
				heap.tree[i] = heap.tree[child]
			} else {
				break
			}
		}
		heap.tree[i] = last
		heap.tree[heap.len] = min
		heap.len--
		return min
	}
	return heap.tree[0]
}
