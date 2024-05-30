package common

import "unicode/utf8"

// splitStringByBytes 将字符串按照指定的字节数进行切割
func SplitStringByBytes(s string, size int) []string {
	var result []string

	for len(s) > 0 {
		// 初始切割点
		l := size
		if l > len(s) {
			l = len(s)
		}

		// 确保不在字符中间切割
		for l > 0 && !utf8.ValidString(s[:l]) {
			l--
		}

		// 如果 l 减到 0，说明 size 太小，无法容纳一个完整的字符
		if l == 0 {
			l = len(s)
			for l > 0 && !utf8.ValidString(s[:l]) {
				l--
			}
		}

		result = append(result, s[:l])
		s = s[l:]
	}

	return result
}
