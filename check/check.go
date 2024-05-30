package check

import (
	"gpt-content-audit/common/config"
	logger "gpt-content-audit/common/loggger"
	"strings"
)

func CheckEnvVariable() {
	if config.AuditChannelType == "" {
		logger.FatalLog("环境变量 AUDIT_CHANNEL_TYPE 未设置")
	} else {
		if strings.ToLower(config.AuditChannelType) == "ali" {
			if config.AliAccessKeyId == "" {
				logger.FatalLog("环境变量 ALI_ACCESS_KEY_ID 未设置")
			}
			if config.AliAccessKeySecret == "" {
				logger.FatalLog("环境变量 ALI_ACCESS_KEY_SECRET 未设置")
			}
			if config.AliEndpoint == "" {
				logger.FatalLog("环境变量 ALI_ENDPOINT 未设置")
			}
			if config.AliLabel == "" {
				logger.FatalLog("环境变量 ALI_LABEL 未设置")
			}
		} else if strings.ToLower(config.AuditChannelType) == "baidu" {
			if config.BaiduApiKey == "" {
				logger.FatalLog("环境变量 BAIDU_API_KEY 未设置")
			}
			if config.BaiduSecretKey == "" {
				logger.FatalLog("环境变量 BAIDU_SECRET_KEY 未设置")
			}
			if config.BaiduLabel == "" {
				logger.FatalLog("环境变量 BAIDU_LABEL 未设置")
			}
		} else {
			logger.FatalLog("不支持的 AUDIT_CHANNEL_TYPE ！")

		}
	}

	if config.BaseUrl == "" {
		logger.FatalLog("环境变量 BASE_URL 未设置")
	}

	if config.Authorization == "" {
		logger.FatalLog("环境变量 AUTHORIZATION 未设置")
	}

	logger.SysLog("Environment variable check passed.")
}
