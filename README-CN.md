# AI Endpoint

## 简介
**简体中文** | [English](https://github.com/ramzeng/ai-endpoint/blob/main/README.md)

一款提供统一接口的 AI 接入点代理层，支持多租户，并具有负载均衡、速率限制和日志记录等能力。

## 具备能力
- [x] Azure API 代理
  - [x] 兼容 OpenAI API 格式
  - [x] API Keys 负载均衡
  - [x] 加权轮询
  - [x] 自适应权重
- [x] 多租户
  - [x] 请求鉴权
  - [x] 配置隔离
- [x] 限流能力
  - [x] 租户级别限流
  - [x] 模型级别限流
- [x] 日志记录能力
  - [x] 多通道写入
  - [x] 文件切割
- [x] 容器化部署

## 支持路由
- POST /v1/chat/completions

## 调用方式
```bash
curl --location 'localhost:8080/v1/chat/completions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer xxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
--data '{
    "stream": true,
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "system",
            "content": "You are a helpful assistant."
        },
        {
            "role": "user",
            "content": "Does Azure OpenAI support customer managed keys?"
        },
        {
            "role": "assistant",
            "content": "Yes, customer managed keys are supported by Azure OpenAI."
        },
        {
            "role": "user",
            "content": "Do other Azure Cognitive Services support this too?"
        }
    ]
}'
```

## 依赖拓展
- [gin](https://github.com/gin-gonic/gin)
- [google/uuid](https://github.com/google/uuid)
- [redis/go-redis](https://github.com/redis/go-redis/v9)
- [spf13/cast](https://github.com/spf13/cast)
- [spf13/viper](https://github.com/spf13/viper)
- [go.uber.org/zap](https://github.com/uber-go/zap)
- [gorm.io/mysql](https://github.com/go-gorm/mysql)
- [gorm.io/gorm](https://github.com/go-gorm/gorm)
- [tidwall/gjson](https://github.com/tidwall/gjson)

## 配置文件
```yaml
logger:
  channels:
    - name: app
      filename: /var/log/app.log
      maxSize: 1
      maxAge: 30
      maxBackups: 10
      compress: false
      level: info
    - name: request
      filename: /var/log/request.log
      maxSize: 1
      maxAge: 30
      maxBackups: 10
      compress: false
      level: info
database:
  addr: root:password@tcp(127.0.0.1:3306)/endpoints?charset=utf8mb4&parseTime=True&loc=Local
redis:
  addr: 127.0.0.1:6379
  password: ""
  db: 0
azure:
  openai:
    models:
      - gpt-3.5-turbo
      - gpt-4
      - gpt-4-32k
    peers:
      - key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
        endpoint: https://xxxxx.openai.azure.com
        weight: 20
        deployments:
          - name: gpt-35-turbo
            model: gpt-3.5-turbo
            version: 2023-03-15-preview
          - name: gpt-4
            model: gpt-4
            version: 2023-03-15-preview
          - name: gpt-4-32k
            model: gpt-4-32k
            version: 2023-03-15-preview
```
