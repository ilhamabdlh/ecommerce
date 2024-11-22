package errors

import "fmt"

type ErrorCode string

const (
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrInternalServer ErrorCode = "INTERNAL_SERVER"
	ErrStockNotEnough ErrorCode = "STOCK_NOT_ENOUGH"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
