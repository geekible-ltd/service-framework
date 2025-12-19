package frameworkdto

import "fmt"

// Response represents a standard API response
// @Description Standard API response structure
type ResponseDTO struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
}

// SuccessResponseDTO represents a successful API response
// @Description Successful API response structure
type SuccessResponseDTO struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
}

// ErrorResponseDTO represents an error API response
// @Description Error API response structure
type ErrorResponseDTO struct {
	Success bool           `json:"success" example:"false"`
	Error   ErrorDetailDTO `json:"error"`
}

// ErrorDetail contains error information
// @Description Error detail structure
type ErrorDetailDTO struct {
	Code    string      `json:"code" example:"ERR_001"`
	Message string      `json:"message" example:"An error occurred"`
	Details interface{} `json:"details,omitempty"`
}

// CreatedResponseDTO represents a 201 Created response
// @Description Created response structure
type CreatedResponseDTO struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message" example:"Resource created successfully"`
}

// ListResponse represents a paginated list response
type ListResponseDTO struct {
	Success    bool           `json:"success"`
	Data       interface{}    `json:"data"`
	Pagination *PaginationDTO `json:"pagination,omitempty"`
}

// Pagination represents pagination metadata
type PaginationDTO struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// AppError represents an application-specific error
type ResponseErrorDTO struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ResponseErrorDTO) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
