package utils

import (
	"bytes"
	"common/utils/gconv"
	"sort"
)

// SortByKeys 对象排序赋值
func SortByKeys(params map[string]interface{}, remove bool, ignore ...string) string {
	var (
		exp     = len(ignore)
		signKey = []string{}
	)

	for k := range params {
		if exp > 0 && InArray(k, ignore) {
			continue
		}
		signKey = append(signKey, k)
	}

	sort.Slice(signKey, func(i, j int) bool {
		return signKey[i] < signKey[j]
	})

	return setKeyValue(signKey, params, remove)
}

func setKeyValue(keys []string, params map[string]interface{}, remove bool) string {
	var signStr bytes.Buffer
	for _, k := range keys {
		// value := fmt.Sprintf("%v", params[k])
		value := gconv.String(params[k])
		if remove && value == "" {
			continue
		}

		if 0 < signStr.Len() {
			signStr.WriteString("&")
		}

		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(value)
	}
	return signStr.String()
}
