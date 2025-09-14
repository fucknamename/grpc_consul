package utils

// 固定键

const ()

var (
	// 签名算法类型
	SIGN_TYPE = map[string]func(data interface{}) string{
		"ASCII_ASC_MD5": func(data interface{}) string { // ASCII 从小到大排 + MD5
			return ""
		},
		"ASCII_ASC_MD5_U": func(data interface{}) string { // ASCII 从小到大排 + MD5 + 大写
			return ""
		},
		"ASCII_ASC_MD5_L": func(data interface{}) string { // ASCII 从小到大排 + MD5 + 小写
			return ""
		},
		"ASCII_ASC_CBC": func(data interface{}) string { // ASCII 从小到大排 + CBC
			return ""
		},
		"ASCII_ASC_CBC_U": func(data interface{}) string { // ASCII 从小到大排 + CBC + 大写
			return ""
		},
		"ASCII_ASC_CBC_L": func(data interface{}) string { // ASCII 从小到大排 + CBC + 小写
			return ""
		},
		"ASCII_ASC_ECB": func(data interface{}) string { // ASCII 从小到大排 + ECB
			return ""
		},
		"ASCII_ASC_ECB_U": func(data interface{}) string { // ASCII 从小到大排 + ECB + 大写
			return ""
		},
		"ASCII_ASC_ECB_L": func(data interface{}) string { // ASCII 从小到大排 + ECB + 小写
			return ""
		},
	}
)
