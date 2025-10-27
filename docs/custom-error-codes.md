# 🔧 自定义错误码系统

## 🎯 功能概述

GoI18n-Gin 库现在支持完全自定义的错误码系统，用户可以选择：

1. **✅ 使用内置错误码** - 预定义的常用错误码
2. **✅ 完全自定义错误码** - 不加载内置错误码，完全自己定义
3. **✅ 混合模式** - 在内置错误码基础上添加自定义错误码

## 📋 配置选项

### ResponseConfig 配置

```go
type ResponseConfig struct {
    LoadBuiltin bool `yaml:"load_builtin" json:"load_builtin"`
    AutoInit     bool `yaml:"auto_init" json:"auto_init"`
}
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `LoadBuiltin` | bool | `true` | 是否加载内置错误码 |
| `AutoInit` | bool | `true` | 是否自动初始化 |

### 使用方式

#### 1. 默认模式（加载内置错误码）
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: true, // 加载内置错误码
        AutoInit:     true, // 自动初始化
    },
}
```

#### 2. 完全自定义模式
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: false, // 不加载内置错误码
        AutoInit:     true,  // 自动初始化
    },
}
```

#### 3. 手动初始化模式
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: false, // 不自动加载
        AutoInit:     false, // 不自动初始化
    },
}

// 手动初始化
i18n.InitWithConfig(config)
response.InitCodes(false) // 不加载内置错误码
```

## 🛠️ 错误码管理 API

### 基础操作

#### 注册单个错误码
```go
response.RegisterCustomCode(1000, "USER_NOT_FOUND", 404)
```

#### 批量注册错误码
```go
codes := []response.CodeDefinition{
    {Code: 1000, Message: "USER_NOT_FOUND", HTTPStatus: 404},
    {Code: 1001, Message: "INVALID_USER_ID", HTTPStatus: 400},
    {Code: 1002, Message: "USER_DISABLED", HTTPStatus: 403},
}

response.BatchRegisterCodes(codes)
```

#### 从映射表加载
```go
messages := map[response.Code]string{
    1000: "USER_NOT_FOUND",
    1001: "INVALID_USER_ID",
}

status := map[response.Code]int{
    1000: 404,
    1001: 400,
}

response.LoadCodesFromMap(messages, status)
```

### 高级操作

#### 设置自定义消息（覆盖内置）
```go
response.SetCustomMessage(response.InvalidParam, "CUSTOM_INVALID_PARAM")
```

#### 设置自定义 HTTP 状态码
```go
response.SetHTTPStatus(response.InvalidParam, 422)
```

#### 注销错误码
```go
response.UnregisterCode(1000)
```

#### 重置整个系统
```go
response.ResetCodes() // 清空所有错误码
```

#### 强制加载内置错误码
```go
response.LoadBuiltinCodesForce() // 会覆盖自定义错误码
```

### 查询操作

#### 获取所有已注册的错误码
```go
registeredCodes := response.GetRegisteredCodes()
fmt.Printf("Total codes: %d\n", len(registeredCodes))
```

#### 获取统计信息
```go
stats := response.GetCodeStats()
fmt.Printf("Total: %d, Client: %d, Server: %d, Custom: %d\n",
    stats["total"], stats["client"], stats["server"], stats["custom"])
```

#### 检查是否已初始化
```go
if response.IsInitialized() {
    fmt.Println("Error code system is initialized")
}
```

## 📚 内置错误码列表

### 成功状态
- `Success (0)` - 成功

### 客户端错误 (1000-1999)
- `InvalidParam (1001)` - 参数错误
- `MissingParam (1002)` - 缺少参数
- `InvalidFormat (1003)` - 格式错误
- `Unauthorized (1004)` - 未授权
- `Forbidden (1005)` - 禁止访问
- `NotFound (1006)` - 资源不存在
- `Conflict (1007)` - 冲突
- `TooManyRequests (1008)` - 请求过多
- `RequestTimeout (1009)` - 请求超时

### 用户相关错误 (1100-1199)
- `UserNotFound (1101)` - 用户不存在
- `UserExists (1102)` - 用户已存在
- `InvalidPassword (1103)` - 密码错误
- `AccountLocked (1104)` - 账户锁定
- `AccountDisabled (1105)` - 账户禁用
- `EmailNotVerified (1106)` - 邮箱未验证
- `PhoneNotVerified (1107)` - 手机未验证

### 认证相关错误 (1200-1299)
- `TokenInvalid (1201)` - Token 无效
- `TokenExpired (1202)` - Token 过期
- `RefreshTokenError (1203)` - 刷新 Token 错误
- `LoginRequired (1204)` - 需要登录
- `PermissionDenied (1205)` - 权限不足
- `SessionExpired (1206)` - 会话过期

### 业务逻辑错误 (1300-1399)
- `BusinessError (1301)` - 业务错误
- `DataConflict (1302)` - 数据冲突
- `OperationFailed (1303)` - 操作失败
- `ResourceExhausted (1304)` - 资源耗尽
- `QuotaExceeded (1305)` - 配额超限
- `RateLimited (1306)` - 频率限制

### 文件相关错误 (1400-1499)
- `FileNotFound (1401)` - 文件不存在
- `FileTooLarge (1402)` - 文件过大
- `FileTypeInvalid (1403)` - 文件类型无效
- `UploadFailed (1404)` - 上传失败
- `DownloadFailed (1405)` - 下载失败
- `StorageExhausted (1406)` - 存储空间不足

### 第三方服务错误 (1500-1599)
- `ThirdPartyError (1501)` - 第三方服务错误
- `ServiceUnavailable (1502)` - 服务不可用
- `ExternalAPIError (1503)` - 外部 API 错误
- `NetworkError (1504)` - 网络错误
- `TimeoutError (1505)` - 超时错误

### 服务器错误 (2000-2999)
- `InternalError (2001)` - 内部错误
- `DatabaseError (2002)` - 数据库错误
- `ServiceError (2003)` - 服务错误
- `ConfigurationError (2004)` - 配置错误
- `DependencyError (2005)` - 依赖错误
- `SystemError (2006)` - 系统错误
- `MaintenanceMode (2007)` - 维护模式

### 未知错误
- `UnknownError (9999)` - 未知错误

## 🎯 使用场景和最佳实践

### 1. 小型项目
**推荐：** 使用内置错误码
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: true,
        AutoInit:     true,
    },
}
```

### 2. 中大型项目
**推荐：** 混合模式
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: true,
        AutoInit:     true,
    },
}

// 在初始化后添加业务相关错误码
addBusinessErrorCodes()
```

### 3. 企业级项目
**推荐：** 完全自定义
```go
config := i18n.Config{
    ResponseConfig: i18n.ResponseConfig{
        LoadBuiltin: false,
        AutoInit:     true,
    },
}

// 完全自定义错误码体系
initializeCustomErrorCodes()
```

## 📊 错误码分类

### 按类型分类
```go
response.GetCategory(code) // 返回 ErrorCategory
// - CategorySuccess
// - CategoryClient
// - CategoryServer
// - CategoryUnknown
```

### 按范围分类
```go
response.IsSuccess(code)     // 是否为成功
response.IsClientError(code)  // 是否为客户端错误 (1000-1999)
response.IsServerError(code)  // 是否为服务器错误 (2000-2999)
response.IsError(code)        // 是否为错误状态
```

## 🔧 示例代码

### 示例1：快速开始（使用内置错误码）
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    // 使用默认配置（自动加载内置错误码）
    if err := i18n.Init(); err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(i18n.Middleware())

    r.GET("/user/:id", func(c *gin.Context) {
        if c.Param("id") == "" {
            response.JSON(c, response.InvalidParam, nil)
            return
        }

        // 业务逻辑...
        response.JSON(c, response.Success, map[string]interface{}{
            "user_id": c.Param("id"),
        })
    })

    r.Run(":8080")
}
```

### 示例2：完全自定义错误码
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    config := i18n.Config{
        ResponseConfig: i18n.ResponseConfig{
            LoadBuiltin: false, // 不加载内置错误码
            AutoInit:     true,
        },
    }

    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    // 注册自定义错误码
    registerCustomCodes()

    r := gin.Default()
    r.Use(i18n.Middleware())

    // 路由定义...
    r.Run(":8080")
}

func registerCustomCodes() {
    codes := []response.CodeDefinition{
        {Code: 10000, Message: "USER_NOT_FOUND", HTTPStatus: 404},
        {Code: 10001, Message: "INVALID_REQUEST", HTTPStatus: 400},
        {Code: 10002, Message: "SERVER_ERROR", HTTPStatus: 500},
    }

    response.BatchRegisterCodes(codes)
}
```

### 示例3：混合模式（内置+自定义）
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/chenguowei/go-i18n"
    "github.com/chenguowei/go-i18n/response"
)

func main() {
    config := i18n.Config{
        ResponseConfig: i18n.ResponseConfig{
            LoadBuiltin: true, // 加载内置错误码
            AutoInit:     true,
        },
    }

    if err := i18n.InitWithConfig(config); err != nil {
        panic(err)
    }

    // 添加业务相关的自定义错误码
    addBusinessCodes()

    r := gin.Default()
    r.Use(i18n.Middleware())

    // 路由定义...
    r.Run(":8080")
}

func addBusinessCodes() {
    // 业务错误码（使用5000-5999范围避免冲突）
    businessCodes := []response.CodeDefinition{
        {Code: 5000, Message: "PRODUCT_OUT_OF_STOCK", HTTPStatus: 422},
        {Code: 5001, Message: "PROMOTION_EXPIRED", HTTPStatus: 410},
        {Code: 5002, Message: "COUPON_ALREADY_USED", HTTPStatus: 409},
    }

    response.BatchRegisterCodes(businessCodes)
}
```

## 🚀 运行示例

```bash
# 运行内置错误码示例
cd examples/quickstart
go run .

# 运行完全自定义错误码示例
cd examples/custom-codes
go run .

# 运行混合模式示例
cd examples/hybrid-codes
go run .
```

## ⚠️ 注意事项

1. **错误码范围**：自定义错误码建议使用 10000+ 或特定业务范围，避免与内置错误码冲突
2. **线程安全**：所有操作都是线程安全的，可以在 goroutine 中安全使用
3. **初始化**：确保在使用错误码之前先初始化系统
4. **覆盖行为**：`LoadBuiltinCodesForce()` 会覆盖已存在的自定义错误码

现在您拥有了完全灵活的错误码管理系统！🎉