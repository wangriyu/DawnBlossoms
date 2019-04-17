package version2

// find a more appropriate hash
func Hash(str []byte) (key int) {
	for _, v := range str {
		key += int(v)
	}
	return key
}
