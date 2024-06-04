package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	"gpt-content-audit/model"
	"io"
	"net/http"
	"strings"
)

var moderationUrl = fmt.Sprintf("%s/v1/moderations", config.OpenaiModerationBaseUrl)

func OpenaiAudit[T model.GetUserContent](t T) (model.AuditResponse, error) {

	labels := strings.Split(config.OpenaiModerationLabel, ",")

	var response model.AuditResponse
	response.Channel = "openai"

	for _, totalContent := range t.GetUserContent() {
		for _, content := range common.SplitStringByBytes(totalContent, config.OpenaiModerationAuditContentLength) {
			request := model.OpenAIModerationRequest{
				Input: content,
			}

			jsonData, err := json.Marshal(request)
			if err != nil {
				return model.AuditResponse{}, err
			}

			req, err := http.NewRequest("POST", moderationUrl, bytes.NewBuffer(jsonData))
			if err != nil {
				return model.AuditResponse{}, err
			}

			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenaiModerationApiKey))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return model.AuditResponse{}, err
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil || resp.StatusCode != http.StatusOK {
				return model.AuditResponse{}, err
			}

			if resp.StatusCode != http.StatusOK {
				return model.AuditResponse{}, fmt.Errorf("request moderations error")
			}

			var openaiResp model.OpenAIModerationResponse
			err = json.Unmarshal(bodyBytes, &openaiResp)
			if err != nil {
				return model.AuditResponse{}, err
			}

			if openaiResp.Results[0].Flagged {
				for _, label := range labels {
					openailabel := strings.ReplaceAll(label, "-", "/")
					if openaiResp.Results[0].Categories[openailabel] {
						resultRes := &model.AuditResultResponse{
							//MessageIndex: i,
							Context:    content,
							Label:      label,
							Suggestion: "PASS",
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
