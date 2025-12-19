package frameworkutils

import (
	"net/http"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
// @Description Standard API response structure
type Response struct {
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
	Success bool        `json:"success" example:"false"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
// @Description Error detail structure
type ErrorDetail struct {
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
type ListResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := err.(*frameworkdto.ResponseErrorDTO); ok {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Error: map[string]interface{}{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			},
		})
		return
	}

	// Default to internal server error for unknown errors
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: map[string]interface{}{
			"code":    frameworkconstants.ErrCodeInternalServer,
			"message": "An unexpected error occurred",
			"details": map[string]interface{}{"error": err.Error()},
		},
	})
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusCreated, data, message)
}

// CreatedResponse sends a 201 Created response
func UpdatedResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusAccepted, data, message)
}

// OKResponse sends a 200 OK response
func OKResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusOK, data, message)
}

// NoContentResponse sends a 204 No Content response
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ListResponseWithPagination sends a paginated list response
func ListResponseWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, ListResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, total int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
