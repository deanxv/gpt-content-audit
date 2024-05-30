package baidu

/*
审核子类型，此字段需参照type主类型字段决定其含义：
当type=11时subType取值含义：
0:百度官方默认违禁词库 default
当type=12时subType取值含义：
0:低质灌水 flood、1:暴恐违禁 terrorism、2:文本色情 porn、3:政治敏感 politics、4:恶意推广 ad、5:低俗辱骂 abuse
当type=13时subType取值含义：
0:自定义文本黑名单 black
当type=14时subType取值含义：
0:自定义文本白名单 white
*/
var baiduAuditConclusionNameToType = map[string]int{
	"PASS":   1,
	"BLOCK":  2,
	"REVIEW": 3,
	"FAILED": 4,
}

// 从字符串到枚举的映射，将在运行时自动填充
var baiduAuditConclusionTypeToName = make(map[int]string)

// 初始化函数，用于生成 stringToEnum 映射
func init() {
	for key, value := range baiduAuditConclusionNameToType {
		baiduAuditConclusionTypeToName[value] = key
	}
}

// 根据枚举获取字符串
func BaiduAuditConclusionGetName(e int) string {
	return baiduAuditConclusionTypeToName[e]
}

// 根据字符串获取枚举
func BaiduAuditConclusionGetEnum(s string) int {
	return baiduAuditConclusionNameToType[s]
}
