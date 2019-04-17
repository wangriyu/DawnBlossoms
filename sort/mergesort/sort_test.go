package sort

import (
	"log"
	"math/rand"
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

func TestMergeSortV1(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	ns = MergeSortV1(ns)
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkMergeSortV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		ns = MergeSortV1(ns)
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}
