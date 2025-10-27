package response

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// NewPagination 创建分页信息
func NewPagination(page, perPage, total int) Pagination {
	totalPages := (total + perPage - 1) / perPage
	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// Paginate 处理分页参数
func Paginate(c *gin.Context) (page, perPage int) {
	// 默认值
	page = 1
	perPage = 10

	// 获取页码
	if p := c.Query("page"); p != "" {
		if parsed, err := parsePage(p); err == nil {
			page = parsed
		}
	}

	// 获取每页数量
	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := parsePerPage(pp); err == nil {
			perPage = parsed
		}
	}

	// 限制最大值
	if perPage > 100 {
		perPage = 100
	}

	return page, perPage
}

// parsePage 解析页码
func parsePage(pageStr string) (int, error) {
	if pageStr == "" {
		return 1, nil
	}

	var page int
	_, err := fmt.Sscanf(pageStr, "%d", &page)
	if err != nil {
		return 1, err
	}

	if page < 1 {
		return 1, fmt.Errorf("page must be greater than 0")
	}

	return page, nil
}

// parsePerPage 解析每页数量
func parsePerPage(perPageStr string) (int, error) {
	if perPageStr == "" {
		return 10, nil
	}

	var perPage int
	_, err := fmt.Sscanf(perPageStr, "%d", &perPage)
	if err != nil {
		return 10, err
	}

	if perPage < 1 {
		return 10, fmt.Errorf("per_page must be greater than 0")
	}

	if perPage > 100 {
		return 100, nil
	}

	return perPage, nil
}

// CreateMeta 创建响应元数据
func CreateMeta(c *gin.Context) *Meta {
	meta := &Meta{
		Timestamp: time.Now(),
	}

	// 获取语言信息
	if lang, exists := c.Get("i18n_language"); exists {
		if str, ok := lang.(string); ok {
			meta.Language = str
		}
	}

	// 获取请求ID
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		meta.RequestID = requestID
	} else {
		meta.RequestID = generateRequestID()
	}

	// 获取追踪ID
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		meta.TraceID = traceID
	}

	return meta
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// BuildQuery 构建查询字符串
func BuildQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for key, value := range params {
		if value != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(parts, "&")
}

// AddPaginationLinks 添加分页链接
func AddPaginationLinks(c *gin.Context, pagination Pagination, basePath string) map[string]string {
	links := make(map[string]string)

	// 构建基础查询参数
	queryParams := make(map[string]string)

	// 保留现有查询参数
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// 第一页
	if pagination.Page > 1 {
		queryParams["page"] = "1"
		links["first"] = fmt.Sprintf("%s?%s", basePath, BuildQuery(queryParams))
	}

	// 上一页
	if pagination.HasPrev {
		queryParams["page"] = fmt.Sprintf("%d", pagination.Page-1)
		links["prev"] = fmt.Sprintf("%s?%s", basePath, BuildQuery(queryParams))
	}

	// 下一页
	if pagination.HasNext {
		queryParams["page"] = fmt.Sprintf("%d", pagination.Page+1)
		links["next"] = fmt.Sprintf("%s?%s", basePath, BuildQuery(queryParams))
	}

	// 最后一页
	if pagination.TotalPages > 0 && pagination.Page < pagination.TotalPages {
		queryParams["page"] = fmt.Sprintf("%d", pagination.TotalPages)
		links["last"] = fmt.Sprintf("%s?%s", basePath, BuildQuery(queryParams))
	}

	return links
}

// WrapData 包装响应数据
func WrapData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
	}
}

// WrapList 包装列表数据
func WrapList(items interface{}, pagination Pagination) map[string]interface{} {
	return map[string]interface{}{
		"items":      items,
		"pagination": pagination,
	}
}

// WrapError 包装错误数据
func WrapError(code Code, message string, details interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"details": details,
	}
}

// ValidatePagination 验证分页参数
func ValidatePagination(page, perPage, maxPerPage int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	if perPage < 1 {
		return fmt.Errorf("per_page must be greater than 0")
	}

	if perPage > maxPerPage {
		return fmt.Errorf("per_page cannot exceed %d", maxPerPage)
	}

	return nil
}

// SanitizeString 清理字符串
func SanitizeString(s string) string {
	// 移除前后空白
	s = strings.TrimSpace(s)

	// 移除多余的空格
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// SanitizeStrings 清理字符串数组
func SanitizeStrings(strs []string) []string {
	var result []string
	for _, s := range strs {
		if cleaned := SanitizeString(s); cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}

// BuildErrorResponse 构建错误响应
func BuildErrorResponse(code Code, err error, details interface{}) map[string]interface{} {
	message := GetMessage(code)
	if err != nil {
		message = err.Error()
	}

	return map[string]interface{}{
		"code":    code,
		"message": message,
		"details": details,
	}
}

// HandlePanic 处理 panic
func HandlePanic(c *gin.Context, recovered interface{}) {
	var err error
	if recoveredErr, ok := recovered.(error); ok {
		err = recoveredErr
	} else {
		err = fmt.Errorf("panic: %v", recovered)
	}

	// 记录 panic
	// 这里可以集成日志系统
	fmt.Printf("Panic recovered: %v\n", err)

	// 返回内部错误响应
	ErrorWithMessage(c, InternalError, "Internal server error")
}

// CORSHeaders 设置 CORS 响应头
func CORSHeaders(c *gin.Context, origins []string) {
	origin := c.GetHeader("Origin")

	// 检查是否允许该来源
	allowed := false
	for _, allowedOrigin := range origins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			allowed = true
			break
		}
	}

	if allowed {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Language")
	c.Header("Access-Control-Max-Age", "86400")
	c.Header("Access-Control-Allow-Credentials", "true")
}

// CacheHeaders 设置缓存响应头
func CacheHeaders(c *gin.Context, maxAge time.Duration, etag string) {
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds()))

	if etag != "" {
		c.Header("ETag", etag)
	}
}

// SecurityHeaders 设置安全响应头
func SecurityHeaders(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
}