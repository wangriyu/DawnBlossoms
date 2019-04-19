package logic

// A valid float looks like "1", "+1.1", "-1.3e10", "-1.3e-2".
// A string may have a prefix that is a valid float.
// Implement a function that returns the valid float prefix.
// Don't need to convert to float, just return the valid prefix.
//
// A few examples are:
// "1.1a" -> "1.1"
// "abc" -> ""
// "-1.1e3.3" -> "-1.1e3"
// "-1.1e." -> "-1.1"
// “1e1” -> "1e1"
func FindFloatPrefix(source string) (result []byte) {
	var t uint
label:
	for _, b := range source {
		v := byte(b)
		switch t {
		case 0:
			if v == '+' || v == '-' {
				t = 1
				result = append(result, v)
			} else if v >= '0' && v <= '9' {
				result = append(result, v)
				t = 2
			} else {
				return
			}
		case 2:
			if v >= '0' && v <= '9' {
				result = append(result, v)
			} else if v == '.' {
				result = append(result, v)
				t = 3
			} else if v == 'e' {
				result = append(result, v)
				t = 5
			} else {
				break label
			}
		case 4:
			if v >= '0' && v <= '9' {
				result = append(result, v)
			} else if v == 'e' {
				result = append(result, v)
				t = 5
			} else {
				break label
			}
		case 5:
			if v == '-' || (v >= '0' && v <= '9') {
				t = 6
				result = append(result, v)
			} else if v == '+' {
				t = 6
			} else {
				break label
			}
		default:
			if v >= '0' && v <= '9' {
				result = append(result, v)
				t++
			} else {
				break label
			}
		}
	}
	for len(result) > 0 && (result[len(result)-1] == '.' || result[len(result)-1] == 'e') {
		result = result[:len(result)-1]
	}
	return result
}
