package http

type PaginationResponse struct {
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
	PagesLeft   int      `json:"pages_left"`
	TotalItems  int      `json:"total_items"`
	ItemsLeft   int      `json:"items_left"`
	Limit       int      `json:"limit"`
	OrderBy     string   `json:"order_by,omitempty"`
	SortOrder   string   `json:"sort_order,omitempty"`
	PageLinks   []string `json:"page_links,omitempty"`
}
