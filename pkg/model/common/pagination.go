package common

import (
	"fmt"
	"math"
	"strconv"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type PaginationQuery struct {
	Page       int    // Page number to retrieve.
	Limit      int    // Maximum number of items per page.
	OrderBy    string // Field to order by.
	SortOrder  string // Sorting direction ("asc" for ascending, "desc" for descending).
	Skip       int    // Number of items to skip for pagination.
	BaseURL    string // Base URL for pagination.
	TotalItems int    // Total number of items.
}

func NewPaginationQuery(page, limit int, orderBy, sortOrder, baseURL string) *PaginationQuery {
	return &PaginationQuery{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		SortOrder: sortOrder,
		Skip:      calculateSkip(page, limit),
		BaseURL:   baseURL,
	}
}

func GetPage(page string) int {
	intPage, stringConversionError := strconv.Atoi(page)
	if validator.IsError(stringConversionError) {
		intPage, _ = strconv.Atoi(constants.DefaultPage)
	}
	if intPage <= 0 {
		intPage, _ = strconv.Atoi(constants.DefaultPage)
	}

	return intPage
}

func GetLimit(limit string) int {
	intLimit, stringConversionError := strconv.Atoi(limit)
	if validator.IsError(stringConversionError) {
		intLimit, _ = strconv.Atoi(constants.DefaultLimit)
	}
	if isLimitInvalid(intLimit) {
		intLimit, _ = strconv.Atoi(constants.DefaultLimit)
	}

	return intLimit
}

func calculateSkip(page, limit int) int {
	if page == 0 {
		return page
	}

	return (page - 1) * limit
}

func isLimitInvalid(data int) bool {
	if data <= 0 || data > constants.MaxItemsPerPage {
		return true
	}

	return false
}

func (paginationQuery *PaginationQuery) SetCorrectPage() {
	if paginationQuery.TotalItems <= paginationQuery.Skip {
		paginationQuery.Page = calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
		paginationQuery.Skip = calculateSkip(paginationQuery.Page, paginationQuery.Limit)
	}
}

type PaginationResponse struct {
	Page       int      // Current page number.
	TotalPages int      // Total number of pages.
	PagesLeft  int      // Number of pages remaining.
	TotalItems int      // Total number of items.
	ItemsLeft  int      // Number of items remaining on the current page.
	Limit      int      // Maximum items per page.
	OrderBy    string   // Field used for ordering.
	SortOrder  string   // Sorting direction ("asc" for ascending, "desc" for descending).
	PageLinks  []string // Array of page links.
	BaseURL    string   // Base URL for pagination.
}

func NewPaginationResponse(paginationQuery *PaginationQuery) *PaginationResponse {
	totalPages := calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
	paginationResponse := PaginationResponse{
		Page:       paginationQuery.Page,
		TotalPages: totalPages,
		PagesLeft:  totalPages - paginationQuery.Page,
		TotalItems: paginationQuery.TotalItems,
		ItemsLeft:  calculateItemsLeft(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit),
		Limit:      paginationQuery.Limit,
		OrderBy:    paginationQuery.OrderBy,
		SortOrder:  paginationQuery.SortOrder,
		BaseURL:    paginationQuery.BaseURL,
	}

	paginationResponse.PageLinks = generatePageLinks(paginationResponse, constants.DefaultAmountOfPages) // You can specify the number of pages to show here.
	return &paginationResponse
}

func calculateTotalPages(totalItems, limit int) int {
	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

func calculateItemsLeft(page, totalItems, limit int) int {
	if (totalItems - (page * limit)) < 0 {
		return 0
	}

	return totalItems - (page * limit)
}

func generatePageLinks(paginationResponse PaginationResponse, amountOfPages int) []string {
	// Preallocate memory for the pageLinks slice based on amountOfPages, adding space for potential first and last page links.
	pageLinks := make([]string, 0, amountOfPages+2)

	// Calculate the range of pages to show before and after the current page.
	startPage := paginationResponse.Page - (amountOfPages / 2)
	if startPage < 1 {
		startPage = 1
	}
	if startPage != 1 {
		pageLinks = append(pageLinks, buildPageLink(paginationResponse, 1))
	}

	// Calculate the end page based on the adjusted start page and amountOfPages.
	endPage := startPage + amountOfPages
	if endPage > paginationResponse.TotalPages {
		// Adjust the end page if it exceeds the total number of pages.
		endPage = paginationResponse.TotalPages
		startPage = endPage - amountOfPages
		if startPage < 1 {
			startPage = 1
		}
	}

	// Append the page links to the list, excluding the current page link.
	for index := startPage; index <= endPage; index++ {
		if index != paginationResponse.Page {
			pageLinks = append(pageLinks, buildPageLink(paginationResponse, index))
		}
	}

	// If the last page is not included in the page links, add the last page link to the pageLinks slice.
	if endPage != paginationResponse.TotalPages {
		pageLinks = append(pageLinks, buildPageLink(paginationResponse, paginationResponse.TotalPages))
	}

	return pageLinks
}

// buildPageLink builds the page link with given page number.
func buildPageLink(paginationResponse PaginationResponse, pageNumber int) string {
	baseURL := paginationResponse.BaseURL

	queryParams := fmt.Sprintf(
		"%s=%d&%s=%d&%s=%s&%s=%s",
		constants.Page,
		pageNumber,
		constants.Limit,
		paginationResponse.Limit,
		constants.OrderBy,
		paginationResponse.OrderBy,
		constants.SortOrder,
		paginationResponse.SortOrder,
	)

	// Combine the base URL and query parameters.
	return fmt.Sprintf("%s?%s", baseURL, queryParams)
}
