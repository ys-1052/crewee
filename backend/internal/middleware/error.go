package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

// APIError represents error information in API responses
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents metadata for API responses (pagination, etc.)
type Meta struct {
	Total  int `json:"total,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ErrorCode represents predefined error codes
type ErrorCode string

const (
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

// newErrorResponse creates an error API response
func newErrorResponse(code ErrorCode, message string, details ...string) APIResponse {
	apiError := &APIError{
		Code:    string(code),
		Message: message,
	}

	if len(details) > 0 {
		apiError.Details = details[0]
	}

	return APIResponse{
		Error: apiError,
	}
}

// ErrorHandler provides centralized error handling
func ErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		resp APIResponse
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		switch code {
		case http.StatusNotFound:
			resp = newErrorResponse(
				ErrCodeNotFound,
				"リソースが見つかりません",
			)
		case http.StatusUnauthorized:
			resp = newErrorResponse(
				ErrCodeUnauthorized,
				"認証が必要です",
			)
		case http.StatusForbidden:
			resp = newErrorResponse(
				ErrCodeForbidden,
				"アクセス権限がありません",
			)
		case http.StatusBadRequest:
			details := ""
			if msg, ok := he.Message.(string); ok {
				details = msg
			}
			resp = newErrorResponse(
				ErrCodeValidation,
				"リクエストが無効です",
				details,
			)
		default:
			resp = newErrorResponse(
				ErrCodeInternal,
				"内部エラーが発生しました",
			)
		}
	} else {
		// Handle other error types
		resp = newErrorResponse(
			ErrCodeInternal,
			"内部エラーが発生しました",
		)
	}

	// Log error for debugging
	c.Logger().Error(err)

	// Send JSON response
	if !c.Response().Committed {
		if err := c.JSON(code, resp); err != nil {
			c.Logger().Error(err)
		}
	}
}
