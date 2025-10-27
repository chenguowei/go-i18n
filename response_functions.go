package i18n

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// 获取错误码对应的消息模板
	messageTemplate := GetMessage(code)

	// 使用全局翻译器翻译消息模板
	var translatedMessage string
	if globalResponseTranslator != nil {
		// 没有模板参数时，直接翻译
		translatedMessage = globalResponseTranslator(c, messageTemplate)
	}

	// 如果翻译失败或未找到翻译，使用原始消息模板
	if translatedMessage == "" {
		translatedMessage = messageTemplate
	}

	response := Response{
		Code:    code,
		Message: translatedMessage,
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
	// 获取错误码对应的消息模板
	messageTemplate := GetMessage(code)

	// 使用全局翻译器翻译消息模板
	var translatedMessage string
	if globalResponseTranslator != nil {
		if templateData != nil {
			// 使用翻译器的模板翻译功能
			translatedMessage = globalResponseTranslator(c, messageTemplate, templateData)
		} else {
			// 没有模板参数时，直接翻译
			translatedMessage = globalResponseTranslator(c, messageTemplate)
		}
	}

	// 如果翻译失败或未找到翻译，使用原始消息模板
	if translatedMessage == "" {
		translatedMessage = messageTemplate
		if templateData != nil {
			// 降级到简单的模板替换
			for key, value := range templateData {
				placeholder := "{{." + key + "}}"
				translatedMessage = strings.ReplaceAll(translatedMessage, placeholder, fmt.Sprintf("%v", value))
			}
		}
	}

	// 构建响应
	meta := &Meta{
		Timestamp: time.Now(),
	}

	// 设置语言信息到元数据
	if l, exists := c.Get("i18n_language"); exists {
		if str, ok := l.(string); ok {
			meta.Language = str
		}
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
		Message: translatedMessage,
		Data:    data,
		Meta:    meta,
	}

	c.JSON(httpStatus, response)
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

// SuccessResponse 成功响应的便捷方法
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