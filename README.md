# GPT-Content-Audit

...

## 功能

...

## 配置

### 环境变量

#### 通用

|     变量参数      |                 变量描述                  | 是否必填 | 
|:-------------:|:-------------------------------------:|:----:|
|  AUDIT_TYPE   |         审核类型[ali:阿里、baidu:百度]         |  Y   |  
| AUTHORIZATION |     鉴权密钥，与接口的API-Key保持一致，多个以`,`分隔     |  Y   |
|   BASE_URL    | 审核通过后的请求地址域名（例如http://api.openai.com） |  Y   |
|    ENABLE     |                审核启用开关                 |  N   |

#### 阿里

|           变量参数           |                                                变量描述                                                | 是否必填 | 
|:------------------------:|:--------------------------------------------------------------------------------------------------:|:----:|
|    ALI_ACCESS_KEY_ID     |                                          阿里云平台AccessKeyId                                          |  Y   |  
|  ALI_ACCESS_KEY_SECRET   |                                        阿里云平台AccessKeySecret                                        |  Y   |
|       ALI_ENDPOINT       |                                           阿里云平台Endpoint                                            |  Y   |
|        ALI_LABEL         | 内容审核类型[spam:垃圾、politics:敏感、abuse:辱骂、terrorism:暴恐、porn:鉴黄、flood:灌水、contraband:违禁、ad:广告] （多个以`,`分隔 ） |  Y   |
| ALI_AUDIT_CONTENT_LENGTH |                                        审核文本切割字节长度[默认:4000]                                         |  N   |




