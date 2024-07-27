package http

type HTTPPaginationResponse struct {
	CurrentPage uint64   `json:"current_page"`         // The current page number being viewed.
	TotalPages  uint64   `json:"total_pages"`          // The total number of pages available.
	PagesLeft   uint64   `json:"pages_left"`           // The number of pages left to view.
	TotalItems  uint64   `json:"total_items"`          // The total number of items available.
	ItemsLeft   uint64   `json:"items_left"`           // The number of items left to view.
	Limit       uint64   `json:"limit"`                // The number of items per page.
	OrderBy     string   `json:"order_by,omitempty"`   // The field by which the items are ordered (optional).
	SortOrder   string   `json:"sort_order,omitempty"` // The order in which the items are sorted, e.g., ascending or descending (optional).
	PageLinks   []string `json:"page_links,omitempty"` // A list of links to other pages (optional).
}

func NewHTTPPaginationResponse(currentPage, totalPages, pagesLeft, itemsLeft, totalItems, limit uint64, orderBy, sortOrder string, pageLinks []string) HTTPPaginationResponse {
	return HTTPPaginationResponse{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		PagesLeft:   pagesLeft,
		TotalItems:  totalItems,
		ItemsLeft:   itemsLeft,
		Limit:       limit,
		OrderBy:     orderBy,
		SortOrder:   sortOrder,
		PageLinks:   pageLinks,
	}
}
