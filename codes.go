package i18n

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Code 错误码类型
type Code int

// Response 统一响应结构
type Response struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta 响应元数据
type Meta struct {
	RequestID  string      `json:"request_id,omitempty"`
	Language   string      `json:"language,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	TraceID    string      `json:"trace_id,omitempty"`
	Version    string      `json:"version,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// ResponseTranslator 响应翻译函数接口
type ResponseTranslator func(c *gin.Context, messageID string, templateData ...map[string]interface{}) string

// 全局响应翻译函数
var globalResponseTranslator ResponseTranslator

// SetResponseTranslator 设置全局响应翻译函数
func SetResponseTranslator(translator ResponseTranslator) {
	globalResponseTranslator = translator
}

// GetResponseTranslator 获取全局响应翻译函数
func GetResponseTranslator() ResponseTranslator {
	return globalResponseTranslator
}

// 全局错误码注册表
var (
	codeMessages    map[Code]string
	httpStatusCodes map[Code]int
	mu              sync.RWMutex
	initialized     bool
)

// 内置错误码定义（可选加载）
const (
	// 成功
	Success Code = 0

	// 公共的错误码 10000-10999
	InvalidParam  Code = 10001 // 参数错误
	InternalError Code = 10002 // 内部错误
)

// 内置错误码映射表
var builtinCodeMessages = map[Code]string{
	Success:       "SUCCESS",
	InvalidParam:  "INVALID_PARAM",
	InternalError: "INTERNAL_ERROR",
}

// 内置 HTTP 状态码映射
var builtinHTTPStatusCodes = map[Code]int{
	Success:       200,
	InvalidParam:  200,
	InternalError: 200,
}

// InitCodes 初始化错误码系统
func InitCodes(loadBuiltin bool) {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return
	}

	// 初始化空的映射表
	codeMessages = make(map[Code]string)
	httpStatusCodes = make(map[Code]int)

	// 加载内置错误码（可选）
	if loadBuiltin {
		for code, message := range builtinCodeMessages {
			codeMessages[code] = message
		}
		for code, status := range builtinHTTPStatusCodes {
			httpStatusCodes[code] = status
		}
	}

	initialized = true
}

// IsInitialized 检查是否已初始化
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return initialized
}

// ResetCodes 重置错误码系统（清空所有注册的错误码）
func ResetCodes() {
	mu.Lock()
	defer mu.Unlock()
	codeMessages = make(map[Code]string)
	httpStatusCodes = make(map[Code]int)
	initialized = false
}

// ensureInitialized 确保错误码系统已初始化（避免锁重入）
func ensureInitialized() {
	if !initialized {
		mu.Lock()
		defer mu.Unlock()

		// 双重检查，防止在等待锁时其他 goroutine 已经初始化
		if !initialized {
			// 初始化空的映射表
			codeMessages = make(map[Code]string)
			httpStatusCodes = make(map[Code]int)

			// 加载内置错误码
			for code, message := range builtinCodeMessages {
				codeMessages[code] = message
			}
			for code, status := range builtinHTTPStatusCodes {
				httpStatusCodes[code] = status
			}

			initialized = true
		}
	}
}

// LoadBuiltinCodes 加载内置错误码（可以重复调用）
func LoadBuiltinCodes() {
	mu.Lock()
	defer mu.Unlock()

	if !initialized {
		codeMessages = make(map[Code]string)
		httpStatusCodes = make(map[Code]int)
		initialized = true
	}

	// 加载内置错误码（不会覆盖已存在的自定义错误码）
	for code, message := range builtinCodeMessages {
		if _, exists := codeMessages[code]; !exists {
			codeMessages[code] = message
		}
	}
	for code, status := range builtinHTTPStatusCodes {
		if _, exists := httpStatusCodes[code]; !exists {
			httpStatusCodes[code] = status
		}
	}
}

// LoadBuiltinCodesForce 强制加载内置错误码（会覆盖已存在的错误码）
func LoadBuiltinCodesForce() {
	mu.Lock()
	defer mu.Unlock()

	if !initialized {
		codeMessages = make(map[Code]string)
		httpStatusCodes = make(map[Code]int)
		initialized = true
	}

	// 强制加载内置错误码（会覆盖已存在的自定义错误码）
	for code, message := range builtinCodeMessages {
		codeMessages[code] = message
	}
	for code, status := range builtinHTTPStatusCodes {
		httpStatusCodes[code] = status
	}
}

// GetMessage 获取错误码对应的消息
func GetMessage(code Code) string {
	ensureInitialized()

	mu.RLock()
	defer mu.RUnlock()

	if message, exists := codeMessages[code]; exists {
		return message
	}
	return "UNKNOWN_ERROR"
}

// GetHTTPStatus 获取错误码对应的 HTTP 状态码
func GetHTTPStatus(code Code) int {
	ensureInitialized()

	mu.RLock()
	defer mu.RUnlock()

	if status, exists := httpStatusCodes[code]; exists {
		return status
	}
	return 200
}

// IsSuccess 判断是否为成功状态
func IsSuccess(code Code) bool {
	return code == Success
}

// IsError 判断是否为错误状态
func IsError(code Code) bool {
	return code != Success
}

// ErrorCategory 错误分类
type ErrorCategory string

const (
	CategorySuccess ErrorCategory = "success"
	CategoryError   ErrorCategory = "error"
)

// GetCategory 获取错误分类
func GetCategory(code Code) ErrorCategory {
	if code == Success {
		return CategorySuccess
	}
	return CategoryError
}

// SetCustomMessage 设置自定义消息
func SetCustomMessage(code Code, message string) {
	ensureInitialized()

	mu.Lock()
	defer mu.Unlock()

	codeMessages[code] = message
}

// SetHTTPStatus 设置自定义 HTTP 状态码
func SetHTTPStatus(code Code, status int) {
	ensureInitialized()

	mu.Lock()
	defer mu.Unlock()

	httpStatusCodes[code] = status
}

// RegisterCustomCode 注册自定义错误码
func RegisterCustomCode(code Code, message string, httpStatus int) {
	ensureInitialized()

	mu.Lock()
	defer mu.Unlock()

	codeMessages[code] = message
	httpStatusCodes[code] = httpStatus
}

// UnregisterCode 注销错误码
func UnregisterCode(code Code) {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		delete(codeMessages, code)
		delete(httpStatusCodes, code)
	}
}

// GetRegisteredCodes 获取所有已注册的错误码
func GetRegisteredCodes() map[Code]string {
	mu.RLock()
	defer mu.RUnlock()

	if !initialized {
		return make(map[Code]string)
	}

	// 返回副本
	result := make(map[Code]string)
	for code, message := range codeMessages {
		result[code] = message
	}
	return result
}

// GetCodeStats 获取错误码统计信息
func GetCodeStats() map[string]int {
	mu.RLock()
	defer mu.RUnlock()

	stats := map[string]int{
		"total":   0,
		"success": 0,
		"error":   0,
		"custom":  0,
	}

	if !initialized {
		return stats
	}

	for code := range codeMessages {
		stats["total"]++

		// 检查是否为内置错误码
		if _, isBuiltin := builtinCodeMessages[code]; isBuiltin {
			if code == Success {
				stats["success"]++
			} else {
				stats["error"]++
			}
		} else {
			stats["custom"]++
		}
	}

	return stats
}

// BatchRegisterCodes 批量注册自定义错误码
type CodeDefinition struct {
	Code       Code
	Message    string
	HTTPStatus int
}

func BatchRegisterCodes(codes []CodeDefinition) {
	ensureInitialized()

	mu.Lock()
	defer mu.Unlock()

	for _, def := range codes {
		codeMessages[def.Code] = def.Message
		httpStatusCodes[def.Code] = def.HTTPStatus
	}
}

// LoadCodesFromMap 从映射表加载错误码
func LoadCodesFromMap(messages map[Code]string, status map[Code]int) {
	ensureInitialized()

	mu.Lock()
	defer mu.Unlock()

	for code, message := range messages {
		codeMessages[code] = message
	}
	for code, httpStatus := range status {
		httpStatusCodes[code] = httpStatus
	}
}
