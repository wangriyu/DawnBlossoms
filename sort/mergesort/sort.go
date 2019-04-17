package sort

func MergeSortV1(source []int) []int {
	if len(source) <= 1 {
		return source
	}

	m := len(source) / 2
	left := MergeSortV1(source[:m])
	right := MergeSortV1(source[m:])

	return mergeV1(left, right)
}

func mergeV1(left, right []int) []int {
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
