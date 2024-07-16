package common

import (
	"fmt"
	"math"
	"strconv"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// PaginationQuery represents the parameters for paginating a list of items.
type PaginationQuery struct {
	Page       int    // Page number to retrieve.
	Limit      int    // Maximum number of items per page.
	OrderBy    string // Field to order by.
	SortOrder  string // Sorting direction ("asc" for ascending, "desc" for descending).
	Skip       int    // Number of items to skip for pagination.
	BaseURL    string // Base URL for pagination.
	TotalItems int    // Total number of items.
}

// NewPaginationQuery creates a new PaginationQuery with the given parameters.
//
// Parameters:
// - page: The page number to retrieve.
// - limit: The maximum number of items per page.
// - orderBy: The field to order by.
// - sortOrder: The sorting direction ("asc" for ascending, "desc" for descending).
// - baseURL: The base URL for pagination.
//
// Returns:
// - A new instance of PaginationQuery with calculated skip value.
func NewPaginationQuery(page, limit int, orderBy, sortOrder, baseURL string) PaginationQuery {
	return PaginationQuery{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		SortOrder: sortOrder,
		Skip:      calculateSkip(page, limit),
		BaseURL:   baseURL,
	}
}

// GetPage maps a string to an integer representing the page number.
//
// Parameters:
// - page: The string representing the page number.
//
// Returns:
// - The integer value of the page number. Defaults to constants.DefaultPage if conversion fails or the value is invalid.
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

// GetLimit maps a string to an integer representing the maximum items per page.
//
// Parameters:
// - limit: The string representing the maximum items per page.
//
// Returns:
// - The integer value of the limit. Defaults to constants.DefaultLimit if conversion fails or the value is invalid.
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

// calculateSkip calculates the number of items to skip based on the current page and limit.
//
// Parameters:
// - page: The current page number.
// - limit: The maximum number of items per page.
//
// Returns:
// - The number of items to skip for pagination.
func calculateSkip(page, limit int) int {
	if page == 0 {
		return page
	}

	return (page - 1) * limit
}

// isLimitInvalid checks if a limit value is valid.
//
// Parameters:
// - data: The limit value to check.
//
// Returns:
// - True if the limit value is invalid, otherwise false.
func isLimitInvalid(data int) bool {
	if data <= 0 || data > constants.MaxItemsPerPage {
		return true
	}

	return false
}

// SetCorrectPage adjusts the PaginationQuery to ensure it's valid, especially when there are not enough items to reach the current page.
//
// Parameters:
// - paginationQuery: The current pagination query.
//
// Returns:
// - The adjusted PaginationQuery with corrected page and skip values.
func SetCorrectPage(paginationQuery PaginationQuery) PaginationQuery {
	if paginationQuery.TotalItems <= paginationQuery.Skip {
		paginationQuery.Page = calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
		paginationQuery.Skip = calculateSkip(paginationQuery.Page, paginationQuery.Limit)
	}

	return paginationQuery
}

// PaginationResponse represents information about the current page, total pages, and more for a paginated list.
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

// NewPaginationResponse creates a new PaginationResponse with the given parameters.
//
// Parameters:
// - paginationQuery: The current pagination query.
//
// Returns:
// - A new instance of PaginationResponse with calculated values and generated page links.
func NewPaginationResponse(paginationQuery PaginationQuery) PaginationResponse {
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
	return paginationResponse
}

// calculateTotalPages calculates the total number of pages based on total items and limit.
//
// Parameters:
// - totalItems: The total number of items.
// - limit: The maximum number of items per page.
//
// Returns:
// - The total number of pages.
func calculateTotalPages(totalItems, limit int) int {
	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

// calculateItemsLeft calculates the number of items remaining on the current page.
//
// Parameters:
// - page: The current page number.
// - totalItems: The total number of items.
// - limit: The maximum number of items per page.
//
// Returns:
// - The number of items remaining on the current page.
func calculateItemsLeft(page, totalItems, limit int) int {
	if (totalItems - (page * limit)) < 0 {
		return 0
	}

	return totalItems - (page * limit)
}

// generatePageLinks creates page links for pagination based on the current page, total pages, and other parameters.
//
// Parameters:
// - paginationResponse: The current pagination response containing pagination details.
// - amountOfPages: The number of pages to show in the pagination links.
//
// Returns:
// - An array of page links.
func generatePageLinks(paginationResponse PaginationResponse, amountOfPages int) []string {
	// Preallocate memory for the pageLinks slice based on amountOfPages, adding space for potential first and last page links.
	pageLinks := make([]string, 0, amountOfPages+2)

	// Calculate the range of pages to show before and after the current page.
	startPage := paginationResponse.Page - (amountOfPages / 2)
	if startPage < 1 {
		startPage = 1
	}

	// If the start page is not the first page, add the first page link to the pageLinks slice.
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
//
// Parameters:
// - paginationResponse: The current pagination response containing pagination details.
// - pageNumber: The page number to build the link for.
//
// Returns:
// - The constructed page link as a string.
func buildPageLink(paginationResponse PaginationResponse, pageNumber int) string {
	baseURL := paginationResponse.BaseURL

	// Construct the query parameters string.
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
