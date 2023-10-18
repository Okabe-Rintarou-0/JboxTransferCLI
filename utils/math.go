package utils

func Max[T int | int64 | float32 | float64](x T, y T) T {
	if x < y {
		return y
	}
	return x
}

func Min[T int | int64 | float32 | float64](x T, y T) T {
	if x > y {
		return y
	}
	return x
}
