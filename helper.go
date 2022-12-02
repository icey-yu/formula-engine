package formula_engine

// IsDigit 是否是数字
func IsDigit(c uint8) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

// IsAlpha 是否是字母
func IsAlpha(c uint8) bool {
	return IsUpAlpha(c) || IsLoAlpha(c)
}

// IsUpAlpha 是否是大写字母
func IsUpAlpha(c uint8) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}
	return false
}

// IsLoAlpha 是否是小写字母
func IsLoAlpha(c uint8) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	return false
}

// InSlice sub是否在s中
func InSlice[T uint8 | string | TT](s []T, sub T) bool {
	for _, em := range s {
		if em == sub {
			return true
		}
	}
	return false
}
