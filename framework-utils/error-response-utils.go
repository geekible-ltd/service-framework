package frameworkutils

import (
	"fmt"
	"net/http"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
)

// NewResponseError creates a new ResponseError
func NewResponseError(code string, message string, statusCode int) *frameworkdto.ResponseErrorDTO {
	return &frameworkdto.ResponseErrorDTO{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// Common errors
func BadRequest(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeBadRequest, message, http.StatusBadRequest)
}

func Unauthorized(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeForbidden, message, http.StatusForbidden)
}

func NotFound(resource string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeNotFound,
		fmt.Sprintf("%s not found", resource),
		http.StatusNotFound,
	)
}

func Conflict(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeConflict, message, http.StatusConflict)
}

func ValidationError(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeValidation, message, http.StatusBadRequest)
}

func InternalServerError(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(frameworkconstants.ErrCodeInternalServer, message, http.StatusInternalServerError)
}

func DatabaseError(err error) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeDatabase,
		"Database operation failed",
		http.StatusInternalServerError,
	)
}

func InvalidInput(field string, reason string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeInvalidInput,
		fmt.Sprintf("Invalid input for field '%s': %s", field, reason),
		http.StatusBadRequest,
	)
}

func MissingHeader(header string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeMissingHeader,
		fmt.Sprintf("Missing required header: %s", header),
		http.StatusBadRequest,
	)
}

func InvalidUUID(field string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeInvalidUUID,
		fmt.Sprintf("Invalid UUID format for field '%s'", field),
		http.StatusBadRequest,
	)
}

func DuplicateEntry(resource string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeDuplicateEntry,
		fmt.Sprintf("%s already exists", resource),
		http.StatusConflict,
	)
}

func ForeignKeyViolation(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeForeignKeyViolation,
		message,
		http.StatusBadRequest,
	)
}

func UnauthorizedError(message string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeUnauthorized,
		message,
		http.StatusUnauthorized,
	)
}

func VersionExistsError(versionNumber string) *frameworkdto.ResponseErrorDTO {
	return NewResponseError(
		frameworkconstants.ErrCodeConflict,
		fmt.Sprintf("Version %s already exists", versionNumber),
		http.StatusConflict,
	)
}
