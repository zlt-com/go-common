package common

// Add 加法
func Add(a, b int) int {
	return a + b
}

// Minus 减法
func Minus(a, b int) int {
	return a - b
}

// Zuijin 整形
// Zuijin 最接近this的数
func Zuijin(this int, arr []int) int {
	min := 0
	if this == arr[0] {
		return arr[0]
	} else if this > arr[0] {
		min = this - arr[0]
	} else if this < arr[0] {
		min = arr[0] - this
	}

	for _, v := range arr {
		if v == this {
			return v
		} else if v > this {
			if min > v-this {
				min = v - this
			}
		} else if v < this {
			if min > this-v {
				min = this - v
			}
		}
	}

	for _, v := range arr {
		if this+min == v {
			return v
		} else if this-min == v {
			return v
		}
	}
	return min
}
