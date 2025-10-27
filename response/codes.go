package response

// Code 错误码类型
type Code int

// 常用错误码定义
const (
	// 成功
	Success Code = 0

	// 客户端错误 1000-1999
	InvalidParam    Code = 1001 // 参数错误
	MissingParam    Code = 1002 // 缺少参数
	InvalidFormat   Code = 1003 // 格式错误
	Unauthorized     Code = 1004 // 未授权
	Forbidden        Code = 1005 // 禁止访问
	NotFound         Code = 1006 // 资源不存在
	Conflict         Code = 1007 // 冲突
	TooManyRequests  Code = 1008 // 请求过多
	RequestTimeout   Code = 1009 // 请求超时

	// 用户相关错误 1100-1199
	UserNotFound     Code = 1101 // 用户不存在
	UserExists       Code = 1102 // 用户已存在
	InvalidPassword  Code = 1103 // 密码错误
	AccountLocked    Code = 1104 // 账户锁定
	AccountDisabled  Code = 1105 // 账户禁用
	EmailNotVerified Code = 1106 // 邮箱未验证
	PhoneNotVerified Code = 1107 // 手机未验证

	// 认证相关错误 1200-1299
	TokenInvalid     Code = 1201 // Token 无效
	TokenExpired      Code = 1202 // Token 过期
	RefreshTokenError Code = 1203 // 刷新 Token 错误
	LoginRequired     Code = 1204 // 需要登录
	PermissionDenied  Code = 1205 // 权限不足
	SessionExpired    Code = 1206 // 会话过期

	// 业务逻辑错误 1300-1399
	BusinessError     Code = 1301 // 业务错误
	DataConflict      Code = 1302 // 数据冲突
	OperationFailed   Code = 1303 // 操作失败
	ResourceExhausted  Code = 1304 // 资源耗尽
	QuotaExceeded     Code = 1305 // 配额超限
	RateLimited       Code = 1306 // 频率限制

	// 文件相关错误 1400-1499
	FileNotFound      Code = 1401 // 文件不存在
	FileTooLarge      Code = 1402 // 文件过大
	FileTypeInvalid   Code = 1403 // 文件类型无效
	UploadFailed      Code = 1404 // 上传失败
	DownloadFailed    Code = 1405 // 下载失败
	StorageExhausted  Code = 1406 // 存储空间不足

	// 第三方服务错误 1500-1599
	ThirdPartyError   Code = 1501 // 第三方服务错误
	ServiceUnavailable Code = 1502 // 服务不可用
	ExternalAPIError  Code = 1503 // 外部 API 错误
	NetworkError      Code = 1504 // 网络错误
	TimeoutError      Code = 1505 // 超时错误

	// 服务器错误 2000-2999
	InternalError     Code = 2001 // 内部错误
	DatabaseError     Code = 2002 // 数据库错误
	ServiceError      Code = 2003 // 服务错误
	ConfigurationError Code = 2004 // 配置错误
	DependencyError   Code = 2005 // 依赖错误
	SystemError       Code = 2006 // 系统错误
	MaintenanceMode   Code = 2007 // 维护模式

	// 未知错误
	UnknownError      Code = 9999 // 未知错误
)

// 错误码映射表
var codeMessages = map[Code]string{
	Success:          "SUCCESS",
	InvalidParam:     "INVALID_PARAM",
	MissingParam:     "MISSING_PARAM",
	InvalidFormat:    "INVALID_FORMAT",
	Unauthorized:     "UNAUTHORIZED",
	Forbidden:        "FORBIDDEN",
	NotFound:         "NOT_FOUND",
	Conflict:         "CONFLICT",
	TooManyRequests:  "TOO_MANY_REQUESTS",
	RequestTimeout:   "REQUEST_TIMEOUT",

	UserNotFound:     "USER_NOT_FOUND",
	UserExists:       "USER_EXISTS",
	InvalidPassword:  "INVALID_PASSWORD",
	AccountLocked:    "ACCOUNT_LOCKED",
	AccountDisabled:  "ACCOUNT_DISABLED",
	EmailNotVerified: "EMAIL_NOT_VERIFIED",
	PhoneNotVerified: "PHONE_NOT_VERIFIED",

	TokenInvalid:     "TOKEN_INVALID",
	TokenExpired:      "TOKEN_EXPIRED",
	RefreshTokenError: "REFRESH_TOKEN_ERROR",
	LoginRequired:     "LOGIN_REQUIRED",
	PermissionDenied:  "PERMISSION_DENIED",
	SessionExpired:    "SESSION_EXPIRED",

	BusinessError:    "BUSINESS_ERROR",
	DataConflict:     "DATA_CONFLICT",
	OperationFailed:  "OPERATION_FAILED",
	ResourceExhausted: "RESOURCE_EXHAUSTED",
	QuotaExceeded:    "QUOTA_EXCEEDED",
	RateLimited:      "RATE_LIMITED",

	FileNotFound:     "FILE_NOT_FOUND",
	FileTooLarge:     "FILE_TOO_LARGE",
	FileTypeInvalid:  "FILE_TYPE_INVALID",
	UploadFailed:     "UPLOAD_FAILED",
	DownloadFailed:   "DOWNLOAD_FAILED",
	StorageExhausted: "STORAGE_EXHAUSTED",

	ThirdPartyError:   "THIRD_PARTY_ERROR",
	ServiceUnavailable: "SERVICE_UNAVAILABLE",
	ExternalAPIError:  "EXTERNAL_API_ERROR",
	NetworkError:      "NETWORK_ERROR",
	TimeoutError:      "TIMEOUT_ERROR",

	InternalError:     "INTERNAL_ERROR",
	DatabaseError:     "DATABASE_ERROR",
	ServiceError:      "SERVICE_ERROR",
	ConfigurationError: "CONFIGURATION_ERROR",
	DependencyError:   "DEPENDENCY_ERROR",
	SystemError:       "SYSTEM_ERROR",
	MaintenanceMode:   "MAINTENANCE_MODE",

	UnknownError:      "UNKNOWN_ERROR",
}

// HTTP 状态码映射
var httpStatusCodes = map[Code]int{
	Success:          200,
	InvalidParam:     400,
	MissingParam:     400,
	InvalidFormat:    400,
	Unauthorized:     401,
	Forbidden:        403,
	NotFound:         404,
	Conflict:         409,
	TooManyRequests:  429,
	RequestTimeout:   408,

	UserNotFound:     404,
	UserExists:       409,
	InvalidPassword:  401,
	AccountLocked:    423,
	AccountDisabled:  403,
	EmailNotVerified: 403,
	PhoneNotVerified: 403,

	TokenInvalid:     401,
	TokenExpired:      401,
	RefreshTokenError: 401,
	LoginRequired:     401,
	PermissionDenied:  403,
	SessionExpired:    401,

	BusinessError:    422,
	DataConflict:     409,
	OperationFailed:  422,
	ResourceExhausted: 429,
	QuotaExceeded:    429,
	RateLimited:      429,

	FileNotFound:     404,
	FileTooLarge:     413,
	FileTypeInvalid:  400,
	UploadFailed:     422,
	DownloadFailed:   500,
	StorageExhausted: 507,

	ThirdPartyError:   502,
	ServiceUnavailable: 503,
	ExternalAPIError:  502,
	NetworkError:      503,
	TimeoutError:      504,

	InternalError:     500,
	DatabaseError:     500,
	ServiceError:      500,
	ConfigurationError: 500,
	DependencyError:   500,
	SystemError:       500,
	MaintenanceMode:   503,

	UnknownError:      500,
}

// GetMessage 获取错误码对应的消息
func GetMessage(code Code) string {
	if message, exists := codeMessages[code]; exists {
		return message
	}
	return codeMessages[UnknownError]
}

// GetHTTPStatus 获取错误码对应的 HTTP 状态码
func GetHTTPStatus(code Code) int {
	if status, exists := httpStatusCodes[code]; exists {
		return status
	}
	return httpStatusCodes[UnknownError]
}

// IsSuccess 判断是否为成功状态
func IsSuccess(code Code) bool {
	return code == Success
}

// IsClientError 判断是否为客户端错误
func IsClientError(code Code) bool {
	return code >= 1000 && code < 2000
}

// IsServerError 判断是否为服务器错误
func IsServerError(code Code) bool {
	return code >= 2000 && code < 3000
}

// IsError 判断是否为错误状态
func IsError(code Code) bool {
	return code != Success
}

// ErrorCategory 错误分类
type ErrorCategory string

const (
	CategorySuccess     ErrorCategory = "success"
	CategoryClient      ErrorCategory = "client_error"
	CategoryServer      ErrorCategory = "server_error"
	CategoryUnknown     ErrorCategory = "unknown"
)

// GetCategory 获取错误分类
func GetCategory(code Code) ErrorCategory {
	if code == Success {
		return CategorySuccess
	}
	if IsClientError(code) {
		return CategoryClient
	}
	if IsServerError(code) {
		return CategoryServer
	}
	return CategoryUnknown
}

// SetCustomMessage 设置自定义消息
func SetCustomMessage(code Code, message string) {
	codeMessages[code] = message
}

// SetHTTPStatus 设置自定义 HTTP 状态码
func SetHTTPStatus(code Code, status int) {
	httpStatusCodes[code] = status
}

// RegisterCustomCode 注册自定义错误码
func RegisterCustomCode(code Code, message string, httpStatus int) {
	codeMessages[code] = message
	httpStatusCodes[code] = httpStatus
}