package http

import "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"

type HTTPPaginationResponse struct {
	CurrentPage int      `json:"current_page"`         // The current page number being viewed.
	TotalPages  int      `json:"total_pages"`          // The total number of pages available.
	PagesLeft   int      `json:"pages_left"`           // The number of pages left to view.
	TotalItems  int      `json:"total_items"`          // The total number of items available.
	ItemsLeft   int      `json:"items_left"`           // The number of items left to view.
	PageStart   int      `json:"page_start"`           // The index of the first item displayed on the current page.
	PageEnd     int      `json:"page_end"`             // The index of the last item displayed on the current page.
	Limit       int      `json:"limit"`                // The number of items per page.
	OrderBy     string   `json:"order_by,omitempty"`   // The field by which the items are ordered (optional).
	SortOrder   string   `json:"sort_order,omitempty"` // The order in which the items are sorted, e.g., ascending or descending (optional).
	PageLinks   []string `json:"page_links,omitempty"` // A list of links to other pages (optional).
}

func NewHTTPPaginationResponse(paginationRespones common.PaginationResponse) HTTPPaginationResponse {
	return HTTPPaginationResponse{
		CurrentPage: paginationRespones.Page,
		TotalPages:  paginationRespones.TotalPages,
		PagesLeft:   paginationRespones.PagesLeft,
		TotalItems:  paginationRespones.TotalItems,
		ItemsLeft:   paginationRespones.ItemsLeft,
		PageStart:   paginationRespones.PageStart,
		PageEnd:     paginationRespones.PageEnd,
		Limit:       paginationRespones.PageEnd,
		OrderBy:     paginationRespones.OrderBy,
		SortOrder:   paginationRespones.SortOrder,
		PageLinks:   paginationRespones.PageLinks,
	}
}
