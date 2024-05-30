package baidu

/*审核子类型，此字段需参照type主类型字段决定其含义：
当type=11时subType取值含义：
0:百度官方默认违禁词库 default
当type=12时subType取值含义：
0:低质灌水 flood、1:暴恐违禁 terrorism、2:文本色情 porn、3:政治敏感 politics、4:恶意推广 ad、5:低俗辱骂 abuse
当type=13时subType取值含义：
0:自定义文本黑名单 black
当type=14时subType取值含义：
0:自定义文本白名单 white
*/

var baiduAuditLabelNameToType = map[string]string{
	"default":   "11-0",
	"flood":     "12-0",
	"terrorism": "12-1",
	"porn":      "12-2",
	"politics":  "12-3",
	"ad":        "12-4",
	"abuse":     "12-5",
	"black":     "13-0",
	"white":     "14-0",
}

// 从字符串到枚举的映射，将在运行时自动填充
var baiduAuditLabelTypeToName = make(map[string]string)

// 初始化函数，用于生成 stringToEnum 映射
func init() {
	for key, value := range baiduAuditLabelNameToType {
		baiduAuditLabelTypeToName[value] = key
	}
}

// 根据枚举获取字符串
func BaiduAuditLabelGetName(e string) string {
	return baiduAuditLabelTypeToName[e]
}

// 根据字符串获取枚举
func BaiduAuditLabelGetEnum(s string) string {
	return baiduAuditLabelNameToType[s]
}
