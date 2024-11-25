package utils

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrWarehouseInactive = errors.New("warehouse is inactive")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden access")
	ErrInternalServer    = errors.New("internal server error")
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}
