package model

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func NewPaginatedResponse(data interface{}, limit, offset, total int) PaginatedResponse {
	return PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}
}
