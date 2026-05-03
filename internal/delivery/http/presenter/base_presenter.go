package presenter

type APIResponse struct {
	Data  interface{} `json:"data"`
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{
		Data: data,
	}
}

func ErrorResponse(err string) APIResponse {
	return APIResponse{
		Error: map[string]string{
			"message": err,
		},
	}
}
