package sort

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	dataLength = 1 << 12
)

var (
	source = make([]int, dataLength)
	max    int
	min    int
)

func init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dataRange := dataLength << 2

	item := r.Intn(dataRange)
	source[0] = item
	min = item
	max = item

	for i := 1; i < dataLength; i++ {
		item := r.Intn(dataRange)
		source[i] = item
		if item < min {
			min = item
		} else if item > max {
			max = item
		}
	}

	log.Printf("Range: %d ~ %d\n", min, max)
}

func TestQuickSortV1(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	QuickSortV1(ns, 0, dataLength-1)
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkQuickSortV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		QuickSortV1(ns, 0, dataLength-1)
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}

func TestQuickSortV2(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	QuickSortV2(ns, 0, dataLength-1)
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkQuickSortV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		QuickSortV2(ns, 0, dataLength-1)
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}

func TestQuickSortV3(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go QuickSortV3(&wg, ns, 0, dataLength-1)
	wg.Wait()
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkQuickSortV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go QuickSortV3(&wg, ns, 0, dataLength-1)
		wg.Wait()
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}

func TestQuickSortV4(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	ch := make(chan struct{}, 1)
	go QuickSortV4(ns, 0, dataLength-1, ch)
	<-ch
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkQuickSortV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		ch := make(chan struct{}, 1)
		go QuickSortV4(ns, 0, dataLength-1, ch)
		<-ch
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}
