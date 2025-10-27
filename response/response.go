package response

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

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
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// JSON 返回 JSON 响应
func JSON(c *gin.Context, code Code, data interface{}) {
	JSONWithMeta(c, code, data, nil)
}

// JSONWithStatus 返回指定 HTTP 状态码的 JSON 响应
func JSONWithStatus(c *gin.Context, code Code, data interface{}, httpStatus int) {
	JSONWithStatusAndMeta(c, code, data, httpStatus, nil)
}

// JSONWithStatusAndMeta 返回指定 HTTP 状态码和元数据的 JSON 响应
func JSONWithStatusAndMeta(c *gin.Context, code Code, data interface{}, httpStatus int, meta *Meta) {
	// 获取语言信息
	lang := "en"
	if l, exists := c.Get("i18n_language"); exists {
		if str, ok := l.(string); ok {
			lang = str
		}
	}

	// 构建响应元数据
	if meta == nil {
		meta = &Meta{}
	}

	meta.Timestamp = time.Now()
	if meta.Language == "" {
		meta.Language = lang
	}

	// 设置请求ID
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		meta.RequestID = requestID
	}

	// 设置追踪ID
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		meta.TraceID = traceID
	}

	response := Response{
		Code:    code,
		Message: GetMessage(code),
		Data:    data,
		Meta:    meta,
	}

	c.JSON(httpStatus, response)
}

// JSONWithMeta 返回带元数据的 JSON 响应
func JSONWithMeta(c *gin.Context, code Code, data interface{}, meta *Meta) {
	JSONWithStatusAndMeta(c, code, data, http.StatusOK, meta)
}

// JSONWithTemplate 支持模板参数的响应
func JSONWithTemplate(c *gin.Context, code Code, data interface{}, templateData map[string]interface{}) {
	JSONWithTemplateAndStatus(c, code, data, templateData, http.StatusOK)
}

// JSONWithTemplateAndStatus 支持模板参数和自定义 HTTP 状态码的响应
func JSONWithTemplateAndStatus(c *gin.Context, code Code, data interface{}, templateData map[string]interface{}, httpStatus int) {
	// 这里可以集成 i18n 模板翻译
	message := GetMessage(code)
	if templateData != nil {
		// 简单的模板替换（实际项目中应该使用更强大的模板引擎）
		for key, value := range templateData {
			placeholder := "{{." + key + "}}"
			message = strings.ReplaceAll(message, placeholder, fmt.Sprintf("%v", value))
		}
	}

	JSONWithStatusAndMeta(c, code, data, httpStatus, &Meta{})
}

// Error 返回错误响应
func Error(c *gin.Context, code Code, message string) {
	JSONWithMeta(c, code, nil, &Meta{
		Timestamp: time.Now(),
	})
}

// ErrorWithMessage 返回自定义错误消息的响应
func ErrorWithMessage(c *gin.Context, code Code, message string) {
	JSONWithStatusAndMeta(c, code, nil, http.StatusOK, &Meta{
		Timestamp: time.Now(),
	})
}

// ErrorWithStatus 返回指定 HTTP 状态码的错误响应
func ErrorWithStatus(c *gin.Context, code Code, httpStatus int) {
	JSONWithStatusAndMeta(c, code, nil, httpStatus, &Meta{
		Timestamp: time.Now(),
	})
}

// ErrorWithMessageAndStatus 返回自定义错误消息和 HTTP 状态码的响应
func ErrorWithMessageAndStatus(c *gin.Context, code Code, message string, httpStatus int) {
	response := Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Meta: &Meta{
			Timestamp: time.Now(),
		},
	}

	c.JSON(httpStatus, response)
}

// Success 成功响应的便捷方法
func SuccessResponse(c *gin.Context, data interface{}) {
	JSON(c, Success, data)
}

// BadRequest 400 错误响应
func BadRequestResponse(c *gin.Context, data interface{}) {
	JSON(c, InvalidParam, data)
}

// Unauthorized 401 错误响应
func UnauthorizedResponse(c *gin.Context, data interface{}) {
	JSON(c, Unauthorized, data)
}

// Forbidden 403 错误响应
func ForbiddenResponse(c *gin.Context, data interface{}) {
	JSON(c, Forbidden, data)
}

// NotFound 404 错误响应
func NotFoundResponse(c *gin.Context, data interface{}) {
	JSON(c, NotFound, data)
}

// InternalServerError 500 错误响应
func InternalServerErrorResponse(c *gin.Context, data interface{}) {
	JSON(c, InternalError, data)
}

// PaginationResponse 分页响应
func PaginationResponse(c *gin.Context, code Code, data interface{}, pagination Pagination) {
	JSONWithMeta(c, code, data, &Meta{
		Pagination: &pagination,
	})
}

// ListResponse 列表响应
func ListResponse(c *gin.Context, code Code, items interface{}, total int, page, perPage int) {
	totalPages := (total + perPage - 1) / perPage
	pagination := Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	JSONWithMeta(c, code, items, &Meta{
		Pagination: &pagination,
	})
}

// SetResponseHeaders 设置响应头
func SetResponseHeaders(c *gin.Context, meta *Meta) {
	if meta != nil {
		if meta.RequestID != "" {
			c.Header("X-Request-ID", meta.RequestID)
		}
		if meta.TraceID != "" {
			c.Header("X-Trace-ID", meta.TraceID)
		}
		if meta.Language != "" {
			c.Header("Content-Language", meta.Language)
		}
	}
}

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	// 根据错误类型返回相应的响应
	if err != nil {
		InternalServerErrorResponse(c, map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		InternalServerErrorResponse(c, nil)
	}
}

// 响应构建器
type Builder struct {
	code         Code
	data         interface{}
	meta         *Meta
	customMessage string
}

// NewBuilder 创建响应构建器
func NewBuilder() *Builder {
	return &Builder{
		code: Success,
		meta: &Meta{},
	}
}

// WithCode 设置响应码
func (b *Builder) WithCode(code Code) *Builder {
	b.code = code
	return b
}

// WithData 设置响应数据
func (b *Builder) WithData(data interface{}) *Builder {
	b.data = data
	return b
}

// WithMeta 设置元数据
func (b *Builder) WithMeta(meta *Meta) *Builder {
	b.meta = meta
	return b
}

// WithLanguage 设置语言
func (b *Builder) WithLanguage(language string) *Builder {
	b.meta.Language = language
	return b
}

// WithRequestID 设置请求ID
func (b *Builder) WithRequestID(requestID string) *Builder {
	b.meta.RequestID = requestID
	return b
}

// WithTraceID 设置追踪ID
func (b *Builder) WithTraceID(traceID string) *Builder {
	b.meta.TraceID = traceID
	return b
}

// WithPagination 设置分页信息
func (b *Builder) WithPagination(pagination Pagination) *Builder {
	b.meta.Pagination = &pagination
	return b
}

// WithCustomMessage 设置自定义消息
func (b *Builder) WithCustomMessage(message string) *Builder {
	b.customMessage = message
	return b
}

// Send 发送响应
func (b *Builder) Send(c *gin.Context) {
	b.meta.Timestamp = time.Now()

	message := b.customMessage
	if message == "" {
		message = GetMessage(b.code)
	}

	response := Response{
		Code:    b.code,
		Message: message,
		Data:    b.data,
		Meta:    b.meta,
	}

	c.JSON(http.StatusOK, response)
}