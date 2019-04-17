package sort

type Data []int

func (dt Data) Len() int {
	return len(dt)
}

func (dt Data) Less(i, j int) bool {
	return dt[i] <= dt[j]
}

func (dt Data) Swap(i, j int) {
	dt[i], dt[j] = dt[j], dt[i]
}

func LessThan(i, j int) bool {
	return i <= j
}
