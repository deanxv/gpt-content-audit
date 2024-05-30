package model

type BaiduComplianceResponse struct {
	ErrorMsg       string          `json:"error_msg"`
	ErrorCode      int             `json:"error_code"`
	Conclusion     string          `json:"conclusion"`
	LogID          int64           `json:"log_id"`
	Data           []BaiduDataItem `json:"data"`
	IsHitMd5       bool            `json:"isHitMd5"`
	ConclusionType int             `json:"conc"`
}

type BaiduDataItem struct {
	Msg            string     `json:"msg"`
	Conclusion     string     `json:"conclusion"`
	Hits           []BaiduHit `json:"hits"`
	SubType        int        `json:"subType"`
	ConclusionType int        `json:"conclusionType"`
	Type           int        `json:"type"`
}

type BaiduHit struct {
	WordHitPositions  []BaiduWordHitPosition `json:"wordHitPositions"`
	Probability       float64                `json:"probability"`
	DatasetName       string                 `json:"datasetName"`
	Words             []string               `json:"words"`
	ModelHitPositions [][]interface{}        `json:"modelHitPositions"`
}

// WordHitPosition structure within Hit
type BaiduWordHitPosition struct {
	Positions [][]int `json:"positions"`
	Label     string  `json:"label"`
}
