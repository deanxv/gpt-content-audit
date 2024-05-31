package model

type QiNiuRequest struct {
	Data   QiNiuTextData  `json:"data"`   // 文本内容
	Params QiNiuParamData `json:"params"` // 请求参数，包括审核类型等
}

type QiNiuTextData struct {
	Text string `json:"text"` // 文本内容
}

type QiNiuParamData struct {
	Scenes []string `json:"scenes"` // 审核类型
}

type QiNiuResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  QiNiuResult `json:"result"`
}

type QiNiuResult struct {
	Suggestion string                      `json:"suggestion"`
	Scenes     map[string]QiNiuSceneDetail `json:"scenes"`
}

type QiNiuSceneDetail struct {
	Suggestion string         `json:"suggestion"`
	Details    []QiNiuDetails `json:"details"`
}

type QiNiuDetails struct {
	Label    string          `json:"label"`
	Score    float64         `json:"score"`
	Contexts []QiNiuContexts `json:"contexts"`
}

type QiNiuContexts struct {
	Context   string          `json:"context"`
	Positions []QiNiuPosition `json:"positions"`
}

type QiNiuPosition struct {
	StartPos int `json:"startPos"`
	EndPos   int `json:"endPos"`
}
