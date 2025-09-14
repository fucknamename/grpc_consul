package utils

// 三元运算
func TernaryT[T comparable](cod bool, r1, r2 T) T {
	if cod {
		return r1
	} else {
		return r2
	}
}
