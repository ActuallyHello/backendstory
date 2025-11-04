package common

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	TotalCount int64 `json:"totalCount"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	HasMore    bool  `json:"hasMore"`
}
