package common

import (
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	_ "github.com/pkoukk/tiktoken-go"
	"strings"
	"unicode/utf8"
)

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

func Obj2Bytes(obj interface{}) ([]byte, error) {
	// 创建一个jsonIter的Encoder
	configCompatibleWithStandardLibrary := jsoniter.ConfigCompatibleWithStandardLibrary
	// 将结构体转换为JSON文本并保持顺序
	bytes, err := configCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GetUUID() string {
	code := uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	return code
}
