package http

type PaginationQuery struct {
	Page  int `query:"page" minimum:"1" default:"1"`
	Limit int `query:"limit" minimum:"1" maximum:"100" default:"20"`
}

func (q PaginationQuery) Offset() int {
	page := q.Page
	if page < 1 {
		page = 1
	}
	limit := q.Limit
	if limit < 1 {
		limit = 20
	}
	return (page - 1) * limit
}

func (q PaginationQuery) SafeLimit() int {
	if q.Limit < 1 {
		return 20
	}
	if q.Limit > 100 {
		return 100
	}
	return q.Limit
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PagedResponse[T any] struct {
	Body struct {
		Data []T            `json:"data"`
		Meta PaginationMeta `json:"meta"`
	} `json:"body"`
}

func NewPagedResponse[T any](items []T, page, limit, total int) *PagedResponse[T] {
	resp := &PagedResponse[T]{}
	resp.Body.Data = items
	resp.Body.Meta = PaginationMeta{Page: page, Limit: limit, Total: total}
	return resp
}
