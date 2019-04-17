package sort

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/psilva261/timsort"
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

func TestTimSort(t *testing.T) {
	ns := make([]int, dataLength)
	copy(ns, source)
	err := timsort.TimSort(Data(ns))
	if err != nil {
		t.Error(err)
	}
	if ns[0] != min || ns[dataLength-1] != max {
		t.Error("unexpected result", ns[0], ns[dataLength-1])
	}
}

func BenchmarkTimSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		err := timsort.TimSort(Data(ns))
		if err != nil {
			b.Error(err)
		}
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}

func BenchmarkTimSortInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ns := make([]int, dataLength)
		copy(ns, source)
		err := timsort.Ints(ns, LessThan)
		if err != nil {
			b.Error(err)
		}
		if ns[0] != min || ns[dataLength-1] != max {
			b.Error("unexpected result", ns[0], ns[dataLength-1])
		}
	}
}
