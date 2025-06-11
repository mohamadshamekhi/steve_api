package response

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewSuccessResponse(message string, data any) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	}
}

func NewErrorResponse(message string, errors any) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
	}
}
