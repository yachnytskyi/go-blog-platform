package http

type PaginationResponse struct {
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	PagesLeft   int    `json:"pages_left"`
	TotalItems  int    `json:"total_items"`
	Limit       int    `json:"limit"`
	OrderBy     string `json:"order_by,omitempty"`
}
