package model

type AuditResponse struct {
	Channel string                 `json:"channel"`
	Results []*AuditResultResponse `json:"results"`
}

type AuditResultResponse struct {
	//MessageIndex int    `json:"messageIndex"`
	Context    string `json:"context"`
	Label      string `json:"label"`
	Suggestion string `json:"suggestion"`
}

func (r AuditResultResponse) ContainsDuplicate(arr []*AuditResultResponse) bool {
	for _, v := range arr {
		// 检查每个元素的Context, Label, Suggestion是否与item相同
		if v.Context == r.Context && v.Label == r.Label && v.Suggestion == r.Suggestion {
			return true
		}
	}

	// 没有发现重复的元素
	return false
}
