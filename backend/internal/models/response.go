package models

// Response представляет общую структуру ответа API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse создает новый успешный ответ
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse создает новый ответ с ошибкой
func NewErrorResponse(err string) Response {
	return Response{
		Success: false,
		Error:   err,
	}
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse представляет ответ с ошибками валидации
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Errors  []ValidationError `json:"errors"`
}

// NewValidationErrorResponse создает новый ответ с ошибками валидации
func NewValidationErrorResponse(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{
		Success: false,
		Errors:  errors,
	}
}
