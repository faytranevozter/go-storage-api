package domain

// BaseResponse ...
type BaseResponse struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Validation map[string]string `json:"validation"`
	Data       interface{}       `json:"data"`
}

// ListResponse ...
type ListResponse struct {
	Limit int64         `json:"limit"`
	List  []interface{} `json:"list"`
	Page  int64         `json:"page"`
	Total int64         `json:"total"`
}
