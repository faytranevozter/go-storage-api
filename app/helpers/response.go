package helpers

import (
	"storage-api/domain"
)

// ErrResp Simple error response
func ErrResp(code int, message string) domain.BaseResponse {
	return domain.BaseResponse{
		Status:     code,
		Validation: make(map[string]string),
		Data:       make(map[string]interface{}),
		Message:    message,
	}
}

// ErrRespVal Error response with validation error
func ErrRespVal(validation map[string]string, message string) domain.BaseResponse {
	return domain.BaseResponse{
		Status:     400,
		Validation: validation,
		Data:       make(map[string]interface{}),
		Message:    message,
	}
}

// SuccessResp Success response
func SuccessResp(i interface{}) domain.BaseResponse {
	return domain.BaseResponse{
		Status:     200,
		Validation: make(map[string]string),
		Data:       i,
		Message:    "success",
	}
}
