package errors

import (
	"fmt"
	"time"
)

// Error 错误接口
type Error interface {
	error
	Code() string
	Message() string
	Details() interface{}
	Timestamp() time.Time
	Cause() error
}

// BaseError 基础错误实现
type BaseError struct {
	code      string
	message   string
	details   interface{}
	timestamp time.Time
	cause     error
}

// New 创建新错误
func New(code, message string) Error {
	return &BaseError{
		code:      code,
		message:   message,
		timestamp: time.Now(),
	}
}

// NewWithDetails 创建带详情的错误
func NewWithDetails(code, message string, details interface{}) Error {
	return &BaseError{
		code:      code,
		message:   message,
		details:   details,
		timestamp: time.Now(),
	}
}

// NewWrap 包装已有错误
func NewWrap(code, message string, cause error) Error {
	return &BaseError{
		code:      code,
		message:   message,
		timestamp: time.Now(),
		cause:     cause,
	}
}

// NewWithCause 创建带原因的错误
func NewWithCause(code, message string, cause error, details interface{}) Error {
	return &BaseError{
		code:      code,
		message:   message,
		details:   details,
		timestamp: time.Now(),
		cause:     cause,
	}
}

// 实现 Error 接口
func (e *BaseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *BaseError) Code() string {
	return e.code
}

func (e *BaseError) Message() string {
	return e.message
}

func (e *BaseError) Details() interface{} {
	return e.details
}

func (e *BaseError) Timestamp() time.Time {
	return e.timestamp
}

func (e *BaseError) Cause() error {
	return e.cause
}

// Is 检查错误是否匹配
func (e *BaseError) Is(target error) bool {
	if t, ok := target.(Error); ok {
		return e.code == t.Code()
	}
	return false
}

// Unwrap 解包错误
func (e *BaseError) Unwrap() error {
	return e.cause
}

// 常用错误码
const (
	// 通用错误
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeInvalidArgument  = "INVALID_ARGUMENT"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodePermissionDenied = "PERMISSION_DENIED"
	ErrCodeUnauthenticated  = "UNAUTHENTICATED"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeCancelled        = "CANCELLED"
	ErrCodeAlreadyExists    = "ALREADY_EXISTS"
	ErrCodeResourceExhausted = "RESOURCE_EXHAUSTED"
	ErrCodeFailedPrecondition = "FAILED_PRECONDITION"
	ErrCodeAborted          = "ABORTED"
	ErrCodeOutOfRange        = "OUT_OF_RANGE"
	ErrCodeUnimplemented     = "UNIMPLEMENTED"
	ErrCodeInternal          = "INTERNAL"
	ErrCodeUnavailable      = "UNAVAILABLE"
	ErrCodeDataLoss          = "DATA_LOSS"

	// 业务错误
	ErrCodeUserNotFound      = "USER_NOT_FOUND"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeAccountLocked     = "ACCOUNT_LOCKED"
	ErrCodeAccountDisabled   = "ACCOUNT_DISABLED"
	ErrCodeEmailNotVerified  = "EMAIL_NOT_VERIFIED"
	ErrCodePhoneNotVerified  = "PHONE_NOT_VERIFIED"
	ErrCodeInvalidToken      = "INVALID_TOKEN"
	ErrCodeTokenExpired      = "TOKEN_EXPIRED"

	// 数据错误
	ErrCodeDatabaseError     = "DATABASE_ERROR"
	ErrCodeConnectionError   = "CONNECTION_ERROR"
	ErrCodeValidationError   = "VALIDATION_ERROR"
	ErrCodeDuplicateKey      = "DUPLICATE_KEY"
	ErrCodeForeignKey        = "FOREIGN_KEY_VIOLATION"
	ErrCodeCheckConstraint   = "CHECK_CONSTRAINT_VIOLATION"

	// 网络错误
	ErrCodeNetworkError      = "NETWORK_ERROR"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeRateLimited       = "RATE_LIMITED"
	ErrCodeRequestTimeout    = "REQUEST_TIMEOUT"

	// 文件错误
	ErrCodeFileNotFound      = "FILE_NOT_FOUND"
	ErrCodeFileTooLarge      = "FILE_TOO_LARGE"
	ErrCodeInvalidFileType   = "INVALID_FILE_TYPE"
	ErrCodeUploadFailed      = "UPLOAD_FAILED"
	ErrCodeStorageError      = "STORAGE_ERROR"

	// 第三方服务错误
	ErrCodeThirdPartyError  = "THIRD_PARTY_ERROR"
	ErrCodeAPIError          = "API_ERROR"
	ErrCodeWebhookError      = "WEBHOOK_ERROR"

	// 配置错误
	ErrCodeConfigError       = "CONFIG_ERROR"
	ErrCodeMissingConfig     = "MISSING_CONFIG"
	ErrCodeInvalidConfig     = "INVALID_CONFIG"
)

// 预定义错误
var (
	ErrInternalError    = New(ErrCodeInternalError, "Internal server error")
	ErrInvalidArgument  = New(ErrCodeInvalidArgument, "Invalid argument")
	ErrNotFound         = New(ErrCodeNotFound, "Resource not found")
	ErrPermissionDenied = New(ErrCodePermissionDenied, "Permission denied")
	ErrUnauthenticated  = New(ErrCodeUnauthenticated, "Unauthenticated")
	ErrTimeout          = New(ErrCodeTimeout, "Request timeout")
	ErrCancelled        = New(ErrCodeCancelled, "Operation cancelled")
	ErrAlreadyExists    = New(ErrCodeAlreadyExists, "Resource already exists")
	ErrResourceExhausted = New(ErrCodeResourceExhausted, "Resource exhausted")
	ErrFailedPrecondition = New(ErrCodeFailedPrecondition, "Failed precondition")
	ErrAborted          = New(ErrCodeAborted, "Operation aborted")
	ErrOutOfRange        = New(ErrCodeOutOfRange, "Argument out of range")
	ErrUnimplemented     = New(ErrCodeUnimplemented, "Not implemented")
	ErrUnavailable      = New(ErrCodeUnavailable, "Service unavailable")
	ErrDataLoss          = New(ErrCodeDataLoss, "Data loss")

	// 业务错误
	ErrUserNotFound      = New(ErrCodeUserNotFound, "User not found")
	ErrInvalidCredentials = New(ErrCodeInvalidCredentials, "Invalid credentials")
	ErrAccountLocked     = New(ErrCodeAccountLocked, "Account is locked")
	ErrAccountDisabled   = New(ErrCodeAccountDisabled, "Account is disabled")
	ErrEmailNotVerified  = New(ErrCodeEmailNotVerified, "Email is not verified")
	ErrPhoneNotVerified  = New(ErrCodePhoneNotVerified, "Phone is not verified")
	ErrInvalidToken      = New(ErrCodeInvalidToken, "Invalid token")
	ErrTokenExpired      = New(ErrCodeTokenExpired, "Token has expired")

	// 数据错误
	ErrDatabaseError     = New(ErrCodeDatabaseError, "Database error")
	ErrConnectionError   = New(ErrCodeConnectionError, "Connection error")
	ErrValidationError   = New(ErrCodeValidationError, "Validation error")
	ErrDuplicateKey      = New(ErrCodeDuplicateKey, "Duplicate key")
	ErrForeignKey        = New(ErrCodeForeignKey, "Foreign key violation")
	ErrCheckConstraint   = New(ErrCodeCheckConstraint, "Check constraint violation")

	// 网络错误
	ErrNetworkError      = New(ErrCodeNetworkError, "Network error")
	ErrServiceUnavailable = New(ErrCodeServiceUnavailable, "Service unavailable")
	ErrRateLimited       = New(ErrCodeRateLimited, "Rate limited")
	ErrRequestTimeout    = New(ErrCodeRequestTimeout, "Request timeout")

	// 文件错误
	ErrFileNotFound      = New(ErrCodeFileNotFound, "File not found")
	ErrFileTooLarge      = New(ErrCodeFileTooLarge, "File too large")
	ErrInvalidFileType   = New(ErrCodeInvalidFileType, "Invalid file type")
	ErrUploadFailed      = New(ErrCodeUploadFailed, "Upload failed")
	ErrStorageError      = New(ErrCodeStorageError, "Storage error")

	// 第三方服务错误
	ErrThirdPartyError  = New(ErrCodeThirdPartyError, "Third party service error")
	ErrAPIError          = New(ErrCodeAPIError, "API error")
	ErrWebhookError      = New(ErrCodeWebhookError, "Webhook error")

	// 配置错误
	ErrConfigError       = New(ErrCodeConfigError, "Configuration error")
	ErrMissingConfig     = New(ErrCodeMissingConfig, "Missing configuration")
	ErrInvalidConfig     = New(ErrCodeInvalidConfig, "Invalid configuration")
)

// Wrap 包装错误
func Wrap(err error, code, message string) Error {
	return NewWrap(code, message, err)
}

// Wrapf 包装错误并格式化消息
func Wrapf(err error, code, format string, args ...interface{}) Error {
	message := fmt.Sprintf(format, args...)
	return NewWrap(code, message, err)
}

// IsCode 检查错误是否为指定代码
func IsCode(err error, code string) bool {
	if e, ok := err.(Error); ok {
		return e.Code() == code
	}
	return false
}

// GetCode 获取错误代码
func GetCode(err error) string {
	if e, ok := err.(Error); ok {
		return e.Code()
	}
	return "UNKNOWN"
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	if e, ok := err.(Error); ok {
		return e.Message()
	}
	return err.Error()
}

// GetDetails 获取错误详情
func GetDetails(err error) interface{} {
	if e, ok := err.(Error); ok {
		return e.Details()
	}
	return nil
}

// GetCause 获取错误原因
func GetCause(err error) error {
	if e, ok := err.(Error); ok {
		return e.Cause()
	}
	return nil
}

// ErrorBuilder 错误构建器
type ErrorBuilder struct {
	code    string
	message string
	details interface{}
	cause   error
}

// NewBuilder 创建错误构建器
func NewBuilder() *ErrorBuilder {
	return &ErrorBuilder{}
}

// Code 设置错误代码
func (b *ErrorBuilder) Code(code string) *ErrorBuilder {
	b.code = code
	return b
}

// Message 设置错误消息
func (b *ErrorBuilder) Message(message string) *ErrorBuilder {
	b.message = message
	return b
}

// Messagef 设置格式化错误消息
func (b *ErrorBuilder) Messagef(format string, args ...interface{}) *ErrorBuilder {
	b.message = fmt.Sprintf(format, args...)
	return b
}

// Details 设置错误详情
func (b *ErrorBuilder) Details(details interface{}) *ErrorBuilder {
	b.details = details
	return b
}

// Cause 设置错误原因
func (b *ErrorBuilder) Cause(cause error) *ErrorBuilder {
	b.cause = cause
	return b
}

// Build 构建错误
func (b *ErrorBuilder) Build() Error {
	if b.cause != nil {
		return NewWithCause(b.code, b.message, b.cause, b.details)
	}
	return NewWithDetails(b.code, b.message, b.details)
}

// 错误类型检查函数
func IsInternalError(err error) bool   { return IsCode(err, ErrCodeInternalError) }
func IsNotFoundError(err error) bool    { return IsCode(err, ErrCodeNotFound) }
func IsPermissionError(err error) bool  { return IsCode(err, ErrCodePermissionDenied) }
func IsTimeoutError(err error) bool      { return IsCode(err, ErrCodeTimeout) }
func IsValidationError(err error) bool   { return IsCode(err, ErrCodeValidationError) }
func IsNetworkError(err error) bool     { return IsCode(err, ErrCodeNetworkError) }
func IsConfigError(err error) bool      { return IsCode(err, ErrCodeConfigError) }
func IsDatabaseError(err error) bool    { return IsCode(err, ErrCodeDatabaseError) }
func IsFileError(err error) bool         { return IsCode(err, ErrCodeFileNotFound) || IsCode(err, ErrCodeFileTooLarge) || IsCode(err, ErrCodeInvalidFileType) }
func IsThirdPartyError(err error) bool   { return IsCode(err, ErrCodeThirdPartyError) }
func IsRateLimitedError(err error) bool  { return IsCode(err, ErrCodeRateLimited) }
func IsServiceUnavailable(err error) bool { return IsCode(err, ErrCodeServiceUnavailable) }