package handlers

import "github.com/ys-1052/crewee/internal/middleware"

// NewSuccessResponse creates a successful API response
func NewSuccessResponse(data interface{}) middleware.APIResponse {
	return middleware.APIResponse{
		Data: data,
	}
}

// NewSuccessResponseWithMeta creates a successful API response with metadata
func NewSuccessResponseWithMeta(data interface{}, meta *middleware.Meta) middleware.APIResponse {
	return middleware.APIResponse{
		Data: data,
		Meta: meta,
	}
}
