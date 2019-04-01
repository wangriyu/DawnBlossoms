package one

// find a more appropriate hash
func Hash(str []byte) (key int) {
	for _, v := range str {
		key += int(v)
	}
	return key % splitFileNum
}
