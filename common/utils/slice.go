package utils

// 从大slice中排除掉小slice
func RemoveSub(a, b []string) []string {
	// 创建一个 map 用于存储需要移除的字符串
	removeMap := make(map[string]struct{}, len(b))
	for _, s := range b {
		removeMap[s] = struct{}{}
	}

	// 遍历数组 a，将不在 map 中的字符串添加到结果数组中
	result := make([]string, 0, len(a))
	for _, s := range a {
		// map 的查询时间复杂度是 O(1)
		if _, ok := removeMap[s]; !ok {
			result = append(result, s)
		}
	}

	return result
}

// 从大slice中排除掉小slice
func RemoveSub2(contain, sub []string) (ret []string) {
	if len(sub) == 0 {
		ret = contain
		return
	}

	for _, v := range contain {
		for _, s := range sub {
			if v != s {
				ret = append(ret, v)
				break
			}
		}
	}

	return
}

// InArray
func InArray[T comparable](needle T, haystack []T) bool {
	for _, item := range haystack {
		if needle == item {
			return true
		}
	}

	return false
}

// 去重
func UniqueArray[T comparable](in []T) (out []T) {
	for _, item := range in {
		if InArray(item, out) {
			continue
		}

		out = append(out, item)
	}
	return
}
