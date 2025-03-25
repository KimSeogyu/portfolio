package db

type CursorBasedPagination struct {
	Cursor *int64 `json:"cursor"`
	Limit  *int64 `json:"limit"`
}

type CursorBasedPaginationResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor *int64 `json:"next_cursor"`
	HasNext    bool   `json:"has_next"`
	Total      int    `json:"total"`
}
