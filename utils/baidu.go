package utils

import (
	"encoding/json"
	"fmt"
	"github.com/deanxv/baidu-golang-sdk/aip/censor"
	"github.com/samber/lo"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	"gpt-content-audit/enum/baidu"
	"gpt-content-audit/model"
	"strings"
)

func BaiduAudit[T model.GetUserContent](t T) (model.AuditResponse, error) {

	// 初始化baidu-client
	client := censor.NewClient(config.BaiduApiKey, config.BaiduSecretKey)
	labels := strings.Split(config.BaiduLabel, ",")

	var response model.AuditResponse
	response.Channel = "baidu"

	// 遍历messages
	for _, totalContent := range t.GetUserContent() {
		for _, content := range common.SplitStringByBytes(totalContent, config.BaiduAuditContentLength) {
			// 发起请求
			resultJson := client.TextCensor(content)
			var resp model.BaiduComplianceResponse
			err := json.Unmarshal([]byte(resultJson), &resp)
			if err != nil {
				return model.AuditResponse{}, err
			}

			if resp.ErrorCode != 0 && resp.ErrorMsg != "" {
				return model.AuditResponse{}, fmt.Errorf(resp.ErrorMsg)
			}

			for _, result := range resp.Data {
				auditLabelName := baidu.BaiduAuditLabelGetName(fmt.Sprintf("%v-%v", result.Type, result.SubType))
				if !lo.Contains(labels, auditLabelName) {
					continue
				}
				conclusionName := baidu.BaiduAuditConclusionGetName(result.ConclusionType)

				for _, hit := range result.Hits {
					// 百度内容审核有bug
					if len(hit.Words) == 0 {
						resultRes := &model.AuditResultResponse{
							//MessageIndex: i,
							Context:    "*",
							Label:      auditLabelName,
							Suggestion: conclusionName,
						}
						if !resultRes.ContainsDuplicate(response.Results) {
							response.Results = append(response.Results, resultRes)
						}
					} else {
						for _, word := range hit.Words {
							resultRes := &model.AuditResultResponse{
								//MessageIndex: i,
								Context:    word,
								Label:      auditLabelName,
								Suggestion: conclusionName,
							}
							if !resultRes.ContainsDuplicate(response.Results) {
								response.Results = append(response.Results, resultRes)
							}
						}
					}
				}
			}
		}

	}

	return response, nil
}
