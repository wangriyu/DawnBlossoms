package main

import (
	"fmt"
	"time"
	"math/rand"
)

func mergeSort(source []int) []int {
	if len(source) <= 1 {
		return source
	}
	m := len(source) / 2
	left := mergeSort(source[:m])
	right := mergeSort(source[m:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	len1 := len(left)
	len2 := len(right)
	tmp := make([]int, 0, len1+len2)
	l := 0
	r := 0
	for l < len1 && r < len2 {
		if left[l] <= right[r] {
			tmp = append(tmp, left[l])
			l++
		} else {
			tmp = append(tmp, right[r])
			r++
		}
	}
	tmp = append(tmp, left[l:]...)
	tmp = append(tmp, right[r:]...)
	return tmp
}

func main() {
	maxLen := 1 << 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := make([]int, maxLen)

	for i := range s {
		s[i] = r.Intn(100)
	}
	fmt.Println("before: ", s)
	s = mergeSort(s)
	fmt.Println("after:  ", s)
}
