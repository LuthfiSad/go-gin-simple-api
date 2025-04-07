package dto

type ResponseError struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Error   map[string]string `json:"errors,omitempty"`
}

type ResponseSuccess struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationMeta struct {
	Page        int   `json:"page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	ItemsOnPage int64 `json:"items_on_page"`
}

type PaginatedResponseData[T any] struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    T              `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}
