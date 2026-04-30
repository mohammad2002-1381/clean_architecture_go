package extensions

type PagedQuery struct {
	PageNumber    int  `json:"page_number" form:"page_number"`
	PageSize      int  `json:"page_size" form:"page_size"`
	DisablePaging bool `json:"disable_paging" form:"disable_paging"`
}

func (q *PagedQuery) Offset() int {
	if q.PageNumber < 1 {
		q.PageNumber = 1
	}
	return (q.PageNumber - 1) * q.PageSize
}

func (q *PagedQuery) Limit() int {
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	return q.PageSize
}

func (q *PagedQuery) Normalize() {
	if q.PageNumber < 1 {
		q.PageNumber = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 10
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

type PagedResult[T any] struct {
	Items      []T  `json:"items"`
	TotalCount int  `json:"total_count"`
	PageNumber int  `json:"page_number"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

func NewPagedResult[T any](items []T, totalCount int, query PagedQuery) PagedResult[T] {
	query.Normalize()
	totalPages := (totalCount + query.PageSize - 1) / query.PageSize

	return PagedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		PageNumber: query.PageNumber,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		HasNext:    query.PageNumber < totalPages,
		HasPrev:    query.PageNumber > 1,
	}
}