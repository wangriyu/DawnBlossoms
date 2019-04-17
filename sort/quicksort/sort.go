package sort

import (
	"math/rand"
	"sync"
	"time"
)

var ra *rand.Rand

func init() {
	ra = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func QuickSortV1(source []int, l, r int) {
	if l < r {
		m := partitionV1(source, l, r)
		QuickSortV1(source, l, m-1)
		QuickSortV1(source, m+1, r)
	}
}

func partitionV1(source []int, l, r int) int {
	k := ra.Intn(r-l) + l
	source[l], source[k] = source[k], source[l]
	pivot := source[l]

	i := l
	for j := l + 1; j <= r; j++ {
		if source[j] < pivot {
			i++
			source[i], source[j] = source[j], source[i]
		}
	}
	source[i], source[l] = source[l], source[i]
	return i
}

func QuickSortV2(source []int, l, r int) {
	if l < r {
		m := partitionV2(source, l, r)
		QuickSortV2(source, l, m-1)
		QuickSortV2(source, m+1, r)
	}
}

func partitionV2(source []int, l, r int) int {
	k := ra.Intn(r-l) + l
	source[l], source[k] = source[k], source[l]
	pivot := source[l]

	i := l
	j := r
	for i < j {
		for i < j && source[j] > pivot {
			j--
		}
		for i < j && source[i] <= pivot {
			i++
		}
		if i < j {
			source[i], source[j] = source[j], source[i]
		}
	}
	source[i], source[l] = source[l], source[i]
	return i
}

func QuickSortV3(parentWg *sync.WaitGroup, source []int, l, r int) {
	defer parentWg.Done()

	if l < r {
		m := partitionV3(source, l, r)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go QuickSortV3(&wg, source, l, m-1)
		go QuickSortV3(&wg, source, m+1, r)
		wg.Wait()
	}
}

func partitionV3(source []int, l, r int) int {
	k := rand.Intn(r-l) + l
	source[l], source[k] = source[k], source[l]
	pivot := source[l]

	i := l
	for j := l + 1; j <= r; j++ {
		if source[j] < pivot {
			i++
			source[i], source[j] = source[j], source[i]
		}
	}
	source[i], source[l] = source[l], source[i]
	return i
}

func QuickSortV4(source []int, l, r int, done chan struct{}) {
	defer func() {
		done <-struct{}{}
	}()

	if l < r {
		m := partitionV4(source, l, r)
		childDone := make(chan struct{}, 2)
		go QuickSortV4(source, l, m-1, childDone)
		go QuickSortV4(source, m+1, r, childDone)
		<-childDone
		<-childDone
	}
}

func partitionV4(source []int, l, r int) int {
	k := rand.Intn(r-l) + l
	source[l], source[k] = source[k], source[l]
	pivot := source[l]

	i := l
	for j := l + 1; j <= r; j++ {
		if source[j] < pivot {
			i++
			source[i], source[j] = source[j], source[i]
		}
	}
	source[i], source[l] = source[l], source[i]
	return i
}
