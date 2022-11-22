package internal

// Reverse[T any]
//
//	@Description: 反序
//	@param arr
func Reverse[T any](arr []T) {
	dLen := len(arr)
	var temp T
	for i := 0; i < dLen/2; i++ {
		temp = arr[i]
		arr[i] = arr[dLen-1-i]
		arr[dLen-1-i] = temp
	}
}

func PowInt64(m int64, n int) int64 {
	total := m
	for i := 0; i < n; i++ {
		total *= m
	}
	return total
}
