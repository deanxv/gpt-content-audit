package config

import (
	"gpt-content-audit/common/env"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ali
var AuditType = os.Getenv("AUDIT_TYPE")
var Enable = env.Int("ENABLE", 1)
var BaseUrl = os.Getenv("BASE_URL")
var Authorization = os.Getenv("AUTHORIZATION")
var Authorizations = strings.Split(os.Getenv("AUTHORIZATION"), ",")

/*
spam：文字垃圾内容识别
politics：文字敏感内容识别
abuse：文字辱骂内容识别
terrorism：文字暴恐内容识别
porn：文字鉴黄内容识别
flood：文字灌水内容识别
contraband：文字违禁内容识别
ad：文字广告内容识别
*/
var AliAccessKeyId = os.Getenv("ALI_ACCESS_KEY_ID")
var AliAccessKeySecret = os.Getenv("ALI_ACCESS_KEY_SECRET")
var AliEndpoint = os.Getenv("ALI_ENDPOINT")
var AliLabel = os.Getenv("ALI_LABEL")
var AliAuditContentLength = env.Int("ALI_AUDIT_CONTENT_LENGTH", 4000) //

/*审核子类型，此字段需参照type主类型字段决定其含义：
当type=11时subType取值含义：
0:百度官方默认违禁词库 default
当type=12时subType取值含义：
0:低质灌水 flood、1:暴恐违禁 terrorism、2:文本色情 porn、3:政治敏感 politics、4:恶意推广 ad、5:低俗辱骂 abuse
当type=13时subType取值含义：
0:自定义文本黑名单 black
当type=14时subType取值含义：
0:自定义文本白名单 white
*/

var BaiduApiKey = os.Getenv("BAIDU_API_KEY")
var BaiduSecretKey = os.Getenv("BAIDU_SECRET_KEY")
var BaiduLabel = os.Getenv("BAIDU_LABEL")
var BaiduAuditContentLength = env.Int("BAIDU_AUDIT_CONTENT_LENGTH", 4000)

var DebugEnabled = strings.ToLower(os.Getenv("DEBUG")) == "true"

var SessionSecret = uuid.New().String()

var RateLimitKeyExpirationDuration = 20 * time.Minute

var (
	RequestRateLimitNum            = env.Int("REQUEST_RATE_LIMIT", 120)
	RequestRateLimitDuration int64 = 1 * 60
)
