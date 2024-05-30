<div align="center">

# GPT-Content-Audit

_聚合阿里、百度等开放平台，提供与`openai`请求格式对齐的内容审核前置服务_

</div>

## 功能

- [x] 支持和`openai`对齐的对话接口(`v1/chat/completions`)用户输入的内容审查。
- [x] 支持和`openai`对齐的`dall-e-3`文生图接口(`v1/images/generations`)用户输入的内容审查。

## 如何使用

1. 配置[环境变量](#环境变量)并[部署](#部署)本项目。
2. 原`openai`请求地址更换为该服务请求地址即可。

## 部署

### 基于 Docker-Compose(All In One) 进行部署

```shell
docker-compose pull && docker-compose up -d
```

#### docker-compose.yml

```docker
version: '3.4'

services:
  gpt-content-audit:
    image: deanxv/gpt-content-audit:latest
    container_name: gpt-content-audit
    restart: always
    ports:
      - "7088:7088"
    volumes:
      - ./data:/app/gpt-content-audit/data
    environment:
      - AUDIT_CHANNEL_TYPE=ali          # 修改为支持的审核渠道类型
      - BASE_URL=https://api.openai.com # 修改为转发后的请求域名或IP:端口
      - AUTHORIZATION=123456            # 修改为转发后的请求地址支持的APIKey
      - TZ=Asia/Shanghai
```

其中`AUDIT_CHANNEL_TYPE`,`BASE_URL`,`AUTHORIZATION`按照自己的需求修改，还需参考[环境变量](#环境变量)配置渠道环境变量。


### 基于 Docker 进行部署

```docker
docker run --name gpt-content-audit -d --restart always \
-p 7088:7088 \
-v $(pwd)/data:/app/gpt-content-audit/data \
-e AUDIT_CHANNEL_TYPE=ali \
-e BASE_URL=https://api.openai.com \
-e AUTHORIZATION=123456 \
-e TZ=Asia/Shanghai \
deanxv/gpt-content-audit
```

如果上面的镜像无法拉取,可以尝试使用 GitHub 的 Docker 镜像,将上面的`deanxv/gpt-content-audit`替换为`ghcr.io/deanxv/gpt-content-audit`即可。

## 配置

### 环境变量

#### 通用

|        变量参数        |                       变量描述                       | 是否必填 | 
|:------------------:|:------------------------------------------------:|:----:|
| AUDIT_CHANNEL_TYPE |             审核渠道类型[ali:阿里、baidu:百度]              |  Y   |  
|      BASE_URL      | 审核通过后的转发接口请求地址域名或IP:端口（例如https://api.openai.com） |  Y   |
|   AUTHORIZATION    |         鉴权密钥，与转发接口的API-Key保持一致，多个以`,`分隔          |  Y   |
|       ENABLE       |             审核启用开关[0:关闭、1:打开]（默认:1）              |  N   |

#### 审核渠道-阿里

|           变量参数           |                                                变量描述                                                | 是否必填 | 
|:------------------------:|:--------------------------------------------------------------------------------------------------:|:----:|
|    ALI_ACCESS_KEY_ID     |                                         阿里云开放平台AccessKeyId                                         |  Y   |  
|  ALI_ACCESS_KEY_SECRET   |                                       阿里云开放平台AccessKeySecret                                       |  Y   |
|       ALI_ENDPOINT       |                                          阿里云开放平台Endpoint                                           |  Y   |
|        ALI_LABEL         | 内容审核类型[spam:垃圾、politics:敏感、abuse:辱骂、terrorism:暴恐、porn:鉴黄、flood:灌水、contraband:违禁、ad:广告] （多个以`,`分隔 ） |  Y   |
| ALI_AUDIT_CONTENT_LENGTH |                                        审核文本切割字节长度[默认:4000]                                         |  N   |

#### 审核渠道-百度

|            变量参数            |                                                  变量描述                                                  | 是否必填 | 
|:--------------------------:|:------------------------------------------------------------------------------------------------------:|:----:|
|       BAIDU_API_KEY        |                                              百度开放平台APIKey                                              |  Y   |  
|      BAIDU_SECRET_KEY      |                                            百度开放平台SecretKey                                             |  Y   |
|        BAIDU_LABEL         | 内容审核类型[default:默认违禁词库、politics:政治敏感、abuse:低俗辱骂、terrorism:暴恐违禁、porn:文本色情、flood:低质灌水、ad:恶意推广]（多个以`,`分隔 ） |  Y   |
| BAIDU_AUDIT_CONTENT_LENGTH |                                          审核文本切割字节长度（默认:4000）                                           |  N   |





