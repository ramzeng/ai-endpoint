# AI Endpoint

## Introduction
**English** | [简体中文](https://github.com/ramzeng/ai-endpoint/blob/main/README-CN.md)

This project is an AI endpoint that provides a unified interface for AI models. It is designed to be multi-tenant, and supports load balancing, rate limiting, and logging capabilities. It is also containerized for easy deployment.

## Capabilities
- [x] Azure API Proxy
  - [x] Compatible with OpenAI API format
  - [x] API Keys load balancing
  - [x] Weighted round-robin
  - [x] Adaptive weight
- [x] Multi-tenant
  - [x] Request authentication
  - [x] Configuration isolation
- [x] Rate limiting capabilities
  - [x] Tenant-level rate limiting
  - [x] Model-level rate limiting
- [x] Logging capabilities
  - [x] Multi-channel writing
  - [x] File splitting
- [x] Containerized deployment

## Supported routes
- POST /v1/chat/completions

## Invocation method
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

## Dependencies
- [gin](https://github.com/gin-gonic/gin)
- [google/uuid](https://github.com/google/uuid)
- [redis/go-redis](https://github.com/redis/go-redis/v9)
- [spf13/cast](https://github.com/spf13/cast)
- [spf13/viper](https://github.com/spf13/viper)
- [go.uber.org/zap](https://github.com/uber-go/zap)
- [gorm.io/mysql](https://github.com/go-gorm/mysql)
- [gorm.io/gorm](https://github.com/go-gorm/gorm)
- [tidwall/gjson](https://github.com/tidwall/gjson)

## Configuration file
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
    backends:
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