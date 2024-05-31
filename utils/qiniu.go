package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	"gpt-content-audit/model"
	"io"
	"net/http"
	"strings"
)

const (
	apiURL      = "https://ai.qiniuapi.com/v3/text/censor"
	contentType = "application/json"
)

func generateQiNiuToken(method, path, rawQuery, host, contentType, bodyStr, accessKey, secretKey string) string {
	var buffer bytes.Buffer
	buffer.WriteString(method + " " + path)
	if rawQuery != "" {
		buffer.WriteString("?" + rawQuery)
	}
	buffer.WriteString("\nHost: " + host)
	if contentType != "" {
		buffer.WriteString("\nContent-Type: " + contentType)
	}
	buffer.WriteString("\n\n")
	if bodyStr != "" && contentType != "" && contentType != "application/octet-stream" {
		buffer.WriteString(bodyStr)
	}

	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write(buffer.Bytes())
	sign := h.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return "Qiniu " + accessKey + ":" + encodedSign
}

func QiNiuAudit[T model.GetUserContent](t T) (model.AuditResponse, error) {

	labels := strings.Split(config.QiNiuLabel, ",")

	var response model.AuditResponse
	response.Channel = "qiniu"

	for _, totalContent := range t.GetUserContent() {
		for _, content := range common.SplitStringByBytes(totalContent, config.QiNiuAuditContentLength) {
			request := model.QiNiuRequest{
				Data: model.QiNiuTextData{
					Text: content,
				},
				Params: model.QiNiuParamData{
					Scenes: []string{"antispam"},
				},
			}

			jsonData, err := json.Marshal(request)
			if err != nil {
				return model.AuditResponse{}, err
			}

			token := generateQiNiuToken("POST", "/v3/text/censor", "", "ai.qiniuapi.com", contentType, string(jsonData), config.QiNiuAccessKey, config.QiNiuSecretKey)

			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
			if err != nil {
				return model.AuditResponse{}, err
			}

			req.Header.Set("Content-Type", contentType)
			req.Header.Set("Authorization", token)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return model.AuditResponse{}, err
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return model.AuditResponse{}, err
			}

			var qiNiuResp model.QiNiuResponse
			err = json.Unmarshal(bodyBytes, &qiNiuResp)
			if err != nil {
				return model.AuditResponse{}, err
			}

			if qiNiuResp.Code != 200 && qiNiuResp.Message != "" {
				return model.AuditResponse{}, fmt.Errorf(qiNiuResp.Message)
			}

			//qiNiuResp.Result.Scenes["antispam"].

			if qiNiuResp.Result.Suggestion != "pass" || qiNiuResp.Result.Scenes["antispam"].Suggestion != "review" {
				if qiNiuResp.Result.Scenes["antispam"].Suggestion != "pass" || qiNiuResp.Result.Scenes["antispam"].Suggestion != "review" {
					for _, detail := range qiNiuResp.Result.Scenes["antispam"].Details {
						if !lo.Contains(labels, detail.Label) {
							continue
						}
						for _, context := range detail.Contexts {
							resultRes := &model.AuditResultResponse{
								//MessageIndex: i,
								Context:    context.Context,
								Label:      detail.Label,
								Suggestion: strings.ToUpper(qiNiuResp.Result.Scenes["antispam"].Suggestion),
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
