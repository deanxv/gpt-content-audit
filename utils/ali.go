package utils

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	imageaudit "github.com/alibabacloud-go/imageaudit-20191230/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	"gpt-content-audit/model"
	"strings"
)

func AliAudit[T model.GetUserContent](t T) (model.AuditResponse, error) {
	// 初始化ali-client
	aliConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AliAccessKeyId),
		AccessKeySecret: tea.String(config.AliAccessKeySecret),
	}
	aliConfig.Endpoint = tea.String(config.AliEndpoint)
	client, err := imageaudit.NewClient(aliConfig)
	if err != nil {
		return model.AuditResponse{}, err
	}

	var tasks []*imageaudit.ScanTextRequestTasks
	// 遍历messages 取出内容
	for _, totalContent := range t.GetUserContent() {
		for _, content := range common.SplitStringByBytes(totalContent, config.AliAuditContentLength) {
			tasks = append(tasks, &imageaudit.ScanTextRequestTasks{
				Content: &content,
			})
		}
	}

	var labels []*imageaudit.ScanTextRequestLabels
	for _, label := range strings.Split(config.AliLabel, ",") {
		labels = append(labels, &imageaudit.ScanTextRequestLabels{
			Label: tea.String(label),
		})
	}

	scanTextRequest := &imageaudit.ScanTextRequest{
		Tasks:  tasks,
		Labels: labels,
	}

	runtime := &util.RuntimeOptions{}
	scanTextResponse, err := client.ScanTextWithOptions(scanTextRequest, runtime)
	if err != nil {
		return model.AuditResponse{}, err
	}

	var response model.AuditResponse
	response.Channel = "ali"

	for _, element := range scanTextResponse.Body.Data.Elements {
		for _, result := range element.Results {
			if *result.Suggestion != "pass" {
				for _, detail := range result.Details {
					for _, context := range detail.Contexts {
						resultRes := &model.AuditResultResponse{
							//MessageIndex: i,
							Context:    *context.Context,
							Label:      *detail.Label,
							Suggestion: strings.ToUpper(*result.Suggestion),
						}
						if !resultRes.ContainsDuplicate(response.Results) {
							response.Results = append(response.Results, resultRes)
						}
					}
				}
			}
		}
	}
	return response, nil
}
