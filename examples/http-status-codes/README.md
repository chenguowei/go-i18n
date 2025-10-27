# HTTP 状态码自定义示例

这个示例展示了如何使用 GoI18n-Gin 库的 HTTP 状态码自定义功能。

## 🎯 功能特性

### 新增的函数

1. **`JSONWithStatus(c, code, data, httpStatus)`**
   - 返回指定 HTTP 状态码的 JSON 响应
   - 适用于需要自定义 HTTP 状态码的场景

2. **`JSONWithStatusAndMeta(c, code, data, httpStatus, meta)`**
   - 返回指定 HTTP 状态码和元数据的 JSON 响应
   - 适用于需要完整控制响应的场景

3. **`ErrorWithStatus(c, code, httpStatus)`**
   - 返回指定 HTTP 状态码的错误响应
   - 适用于错误场景的状态码自定义

4. **`ErrorWithMessageAndStatus(c, code, message, httpStatus)`**
   - 返回自定义错误消息和 HTTP 状态码的响应
   - 适用于需要自定义错误消息的场景

5. **`JSONWithTemplateAndStatus(c, code, data, templateData, httpStatus)`**
   - 返回支持模板参数和自定义 HTTP 状态码的响应
   - 适用于需要动态生成消息模板的场景
   - **支持真正的多语言翻译功能**

## 🌍 多语言翻译功能

### 翻译机制
- 使用内置的 i18n 翻译系统自动翻译错误码消息
- 支持模板参数的多语言翻译
- 通过 `Accept-Language` 头自动检测语言偏好
- 支持多种语言文件格式（JSON、YAML、TOML）

### 使用方式
```go
// 自动根据请求语言翻译错误消息
i18n.JSONWithTemplateAndStatus(c, i18n.Success,
    data, templateData, http.StatusOK)

// 英文请求返回英文消息
curl -H "Accept-Language: en" /api/endpoint

// 中文请求返回中文消息
curl -H "Accept-Language: zh-CN" /api/endpoint
```

## 🚀 运行示例

```bash
cd examples/http-status-codes
go run .
```

## 📋 API 端点说明

### 基础示例

| 方法 | 端点 | 说明 | HTTP 状态码 |
|------|------|------|-------------|
| GET | `/default` | 默认状态码 (200) | 200 |
| POST | `/created` | 资源创建成功 | 201 |
| GET | `/bad-request` | 错误请求 | 400 |
| GET | `/unprocessable` | 无法处理的实体 | 422 |
| PUT | `/accepted` | 请求已接受 | 202 |
| GET | `/template` | 模板参数响应 | 201 |
| GET | `/template/i18n` | 多语言模板响应 | 200 |
| GET | `/template/error` | 多语言错误模板 | 400 |

### RESTful API 示例

| 方法 | 端点 | 说明 | 成功状态码 | 错误状态码 |
|------|------|------|-------------|-------------|
| GET | `/api/v1/users` | 获取用户列表 | 200 | - |
| POST | `/api/v1/users` | 创建用户 | 201 | 400 |
| GET | `/api/v1/users/:id` | 获取特定用户 | 200 | 404 |
| PUT | `/api/v1/users/:id` | 更新用户 | 200 | 404 |
| DELETE | `/api/v1/users/:id` | 删除用户 | 204 | 404 |

### 场景说明

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/scenarios` | 不同业务场景的状态码说明 |

## 💡 使用示例

### 1. 基础用法（支持多语言）

```go
// 默认状态码 (200) - 自动翻译消息
i18n.JSON(c, i18n.Success, data)

// 自定义状态码 (201) - 自动翻译消息
i18n.JSONWithStatus(c, i18n.Success, data, http.StatusCreated)
```

### 2. 多语言响应机制

所有 JSON 响应函数现在都支持自动多语言翻译：

```go
// 英文请求
curl -H "Accept-Language: en" /api/endpoint
// 返回：{"code":0,"message":"Operation successful"}

// 中文请求
curl -H "Accept-Language: zh-CN" /api/endpoint
// 返回：{"code":0,"message":"操作成功"}

// 错误消息也会自动翻译
curl -H "Accept-Language: zh-CN" /api/bad-request
// 返回：{"code":1001,"message":"参数错误"}
```

### 3. 错误响应

```go
// 默认错误状态码 (200)
i18n.Error(c, i18n.InvalidParam)

// 自定义错误状态码 (400)
i18n.ErrorWithStatus(c, i18n.InvalidParam, http.StatusBadRequest)

// 自定义错误消息和状态码
i18n.ErrorWithMessageAndStatus(c, i18n.InvalidParam,
    "参数验证失败", http.StatusUnprocessableEntity)
```

### 3. 带元数据的响应

```go
meta := &i18n.Meta{
    RequestID: "req-123",
    Language:  "zh-CN",
    Version:   "v1.0",
}

i18n.JSONWithStatusAndMeta(c, i18n.Success,
    data, http.StatusCreated, meta)
```

### 4. 模板参数响应

```go
templateData := map[string]interface{}{
    "ResourceName": "用户",
    "ResourceID":   "12345",
    "Action":       "创建",
    "Timestamp":    time.Now().Format("2006-01-02 15:04:05"),
}

i18n.JSONWithTemplateAndStatus(c, i18n.Success,
    data, templateData, http.StatusCreated)
```

## 🔗 相关文档

- [自定义错误码系统](../../docs/custom-error-codes.md)
- [响应码系统演进](../../docs/response-codes-evolution.md)
- [快速开始指南](../../docs/quickstart-guide.md)

## 🎯 最佳实践

1. **RESTful API**: 使用标准 HTTP 状态码
   - 200: 成功获取资源
   - 201: 资源创建成功
   - 204: 资源删除成功
   - 400: 客户端请求错误
   - 401: 未授权
   - 403: 禁止访问
   - 404: 资源不存在
   - 500: 服务器内部错误

2. **业务错误码**: 使用业务错误码 + 合适的 HTTP 状态码
   - `i18n.UserNotFound` + `404`
   - `i18n.InvalidParam` + `400`
   - `i18n.Unauthorized` + `401`

3. **一致性**: 在整个项目中保持状态码使用的一致性