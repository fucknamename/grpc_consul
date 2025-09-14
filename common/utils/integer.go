package utils

import "golang.org/x/exp/constraints"

// Max 返回两个数中的较大者
func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 返回两个数中的较小者
func Min[T constraints.Integer](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// note: constraints.Ordered 约束浮点数
