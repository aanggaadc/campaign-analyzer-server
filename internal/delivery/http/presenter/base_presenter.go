package presenter

type APIResponse struct {
	Data  interface{} `json:"data"`
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func PaginationMeta(page, limit, total int) Meta {
	totalPages := (total + limit - 1) / limit

	return Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

func Success(data interface{}) APIResponse {
	return APIResponse{
		Data: data,
	}
}

func SuccessWithMeta(data interface{}, meta interface{}) APIResponse {
	return APIResponse{
		Data: data,
		Meta: meta,
	}
}

func ErrorResponse(err string) APIResponse {
	return APIResponse{
		Error: map[string]string{
			"message": err,
		},
	}
}
