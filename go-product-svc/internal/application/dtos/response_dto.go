package dtos

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, err interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

func NewValidationErrorResponse(errors []ValidationError) APIResponse {
	return APIResponse{
		Success: false,
		Message: "Validation failed",
		Error: ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "One or more validation errors occurred",
			Details: errors,
		},
	}
}
