package domain

type PaginationParams struct {
	PageNumber    int    `json:"page_number" form:"page_number"`
	PageSize      int    `json:"page_size" form:"page_size"`
	Sort          string `json:"sort" form:"sort"`
	Desc          bool   `json:"desc" form:"desc"`
	DisablePaging bool   `json:"disable_paging" form:"disable_paging"`
}

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total"`
	PageNumber int   `json:"page"`
	PageSize   int   `json:"limit"`
}
