package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	logger "gpt-content-audit/common/loggger"
	"gpt-content-audit/middleware"
	"gpt-content-audit/model"
	"gpt-content-audit/utils"
	"io"
	"net/http"
	"strings"
	"time"
)

func ChatForOpenAI(c *gin.Context) {

	// 读取并缓存请求体
	var bodyBytes bytes.Buffer
	_, err := io.Copy(&bodyBytes, c.Request.Body)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "INVALID_REQUEST",
			},
		})
		return
	}

	// 重设请求体，以便解析JSON
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	// 解析请求
	var request model.OpenAIChatCompletionRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "INVALID_REQUEST",
			},
		})
		return
	}

	// 重设请求体，以便转发
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	// 判断审核渠道
	response := model.AuditResponse{}

	if strings.ToLower(config.AuditChannelType) == "ali" {
		response, err = utils.AliAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: err.Error(),
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "baidu" {
		response, err = utils.BaiduAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: err.Error(),
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "qiniu" {
		response, err = utils.QiNiuAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: err.Error(),
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "openai" {
		response, err = utils.OpenaiAudit(c, request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: err.Error(),
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Unknown audit channel",
				Type:    "request_error",
				Code:    "AUDIT_CHANNEL_ERROR",
			},
		})
		return
	}

	if response.Results == nil || len(response.Results) == 0 {
		middleware.ForwardTo(c, config.BaseUrl)
	} else {
		if config.CustomAuditResult == "" {
			errMsg := ""
			for _, re := range response.Results {
				errMsg += fmt.Sprintf("[%s:%s]", re.Label, re.Context)
			}

			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: fmt.Sprintf("Sensitive information detected:%s", errMsg),
					Type:    "request_error",
					Code:    "AUDIT_RESULT",
				},
			})
			return
		} else {
			res, err := buildRes(request)
			if err != nil {
				logger.Errorf(c.Request.Context(), err.Error())
				c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
					OpenAIError: model.OpenAIError{
						Message: err.Error(),
						Type:    "request_error",
						Code:    "AUDIT_CHANNEL_ERROR",
					},
				})
				return
			}

			if request.Stream {
				c.Stream(func(w io.Writer) bool {
					resBytes, _ := common.Obj2Bytes(res)
					c.SSEvent("", " "+string(resBytes))
					c.SSEvent("", " [DONE]")
					return false
				})
			} else {
				res.Object = "chat.completion"
				c.JSON(http.StatusOK, res)
				return
			}
		}
	}
}

func buildRes(request model.OpenAIChatCompletionRequest) (*model.OpenAIChatCompletionResponse, error) {
	jsonBytes, err := json.Marshal(request.Messages)
	if err != nil {
		return nil, err
	}
	// 将字节切片转换为字符串并打印
	jsonString := string(jsonBytes)

	promptTokens := common.CountTokens(jsonString)
	completionTokens := common.CountTokens(config.CustomAuditResult)
	res := &model.OpenAIChatCompletionResponse{
		ID:      common.GetUUID(),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   request.Model,
		Choices: []model.OpenAIChoice{
			{
				Index: 0,
				Message: model.OpenAIMessage{
					Role:    "assistant",
					Content: config.CustomAuditResult,
				},
				FinishReason: "stop",
				Delta: model.OpenAIDelta{
					Content: config.CustomAuditResult,
				},
			},
		},

		Usage: model.OpenAIUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
	return res, nil
}

func ImagesForOpenAI(c *gin.Context) {
	// 读取并缓存请求体
	var bodyBytes bytes.Buffer
	_, err := io.Copy(&bodyBytes, c.Request.Body)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "INVALID_REQUEST",
			},
		})
		return
	}

	// 重设请求体，以便解析JSON
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	// 解析请求
	var request model.OpenAIImagesGenerationRequest
	err = json.NewDecoder(c.Request.Body).Decode(&request)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "INVALID_REQUEST",
			},
		})
		return
	}

	// 重设请求体，以便转发
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes.Bytes()))

	// 判断审核渠道

	response := model.AuditResponse{}

	if strings.ToLower(config.AuditChannelType) == "ali" {
		response, err = utils.AliAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: "Unknown audit channel",
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "baidu" {
		response, err = utils.BaiduAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: "Unknown audit channel",
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "qiniu" {
		response, err = utils.QiNiuAudit(request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: "Unknown audit channel",
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else if strings.ToLower(config.AuditChannelType) == "openai" {
		response, err = utils.OpenaiAudit(c, request)
		if err != nil {
			logger.Errorf(c.Request.Context(), err.Error())
			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: err.Error(),
					Type:    "request_error",
					Code:    "AUDIT_CHANNEL_ERROR",
				},
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Unknown audit channel",
				Type:    "request_error",
				Code:    "AUDIT_CHANNEL_ERROR",
			},
		})
		return
	}

	if response.Results == nil || len(response.Results) == 0 {
		middleware.ForwardTo(c, config.BaseUrl)
	} else {
		if config.CustomAuditResult == "" {
			errMsg := ""
			for _, re := range response.Results {
				errMsg += fmt.Sprintf("[%s:%s]", re.Label, re.Context)
			}

			c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
				OpenAIError: model.OpenAIError{
					Message: fmt.Sprintf("Sensitive information detected:%s", errMsg),
					Type:    "request_error",
					Code:    "AUDIT_RESULT",
				},
			})
			return
		} else {
			resData := &model.OpenAIImagesGenerationDataResponse{
				RevisedPrompt: config.CustomAuditResult,
			}
			res := model.OpenAIImagesGenerationResponse{
				Created: time.Now().Unix(),
				Data:    []*model.OpenAIImagesGenerationDataResponse{resData},
			}
			c.JSON(http.StatusOK, res)
			return
		}
	}
}
