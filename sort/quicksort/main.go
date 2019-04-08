package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSort(source []int, l, r int) {
	if l < r {
		m := partition(source, l, r)
		quickSort(source, l, m-1)
		quickSort(source, m+1, r)
	}
}

func partition(source []int, l, r int) int {
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
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

func main() {
	maxLen := 1<<5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]int, maxLen)

	for i := range s {
		s[i] = r.Intn(1000) - 10
	}
	fmt.Println("before: ", s)
	quickSort(s, 0, len(s)-1)
	fmt.Println("after:  ", s)
}
