package common

import (
	"fmt"
	"math"
	"strconv"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type PaginationQuery struct {
	Page       uint64 // Page number to retrieve.
	Limit      uint64 // Maximum number of items per page.
	OrderBy    string // Field to order by.
	SortOrder  string // Sorting direction ("asc" for ascending, "desc" for descending).
	Skip       uint64 // Number of items to skip for pagination.
	BaseURL    string // Base URL for pagination.
	TotalItems uint64 // Total number of items.
}

func NewPaginationQuery(page, limit uint64, orderBy, sortOrder, baseURL string) PaginationQuery {
	return PaginationQuery{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		SortOrder: getSortOrder(sortOrder),
		Skip:      calculateSkip(page, limit),
		BaseURL:   baseURL,
	}
}

func GetPage(page string) uint64 {
	uintPage, stringConversionError := strconv.ParseUint(page, 0, 0)
	if validator.IsError(stringConversionError) || uintPage == 0 {
		uintPage, _ = strconv.ParseUint(constants.DefaultPage, 0, 0)
	}

	return uint64(uintPage)
}

func GetLimit(limit string) uint64 {
	uintLimit, stringConversionError := strconv.ParseUint(limit, 0, 0)
	if validator.IsError(stringConversionError) {
		uintLimit, _ = strconv.ParseUint(constants.DefaultLimit, 0, 0)
	}
	if isLimitInvalid(uintLimit) {
		uintLimit, _ = strconv.ParseUint(constants.DefaultLimit, 0, 0)
	}

	return uintLimit
}

func getSortOrder(sortOrder string) string {
	if sortOrder != constants.SortAscend && sortOrder != constants.SortDescend {
		return constants.DefaultSortOrder
	}

	return sortOrder
}

// calculateSkip calculates the number of items to skip for pagination.
func calculateSkip(page, limit uint64) uint64 {
	return (page - 1) * limit
}

func isLimitInvalid(data uint64) bool {
	if data > constants.MaxItemsPerPage {
		return true
	}

	return false
}

func SetCorrectPage(paginationQuery PaginationQuery) PaginationQuery {
	if paginationQuery.TotalItems <= paginationQuery.Skip {
		paginationQuery.Page = calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
		paginationQuery.Skip = calculateSkip(paginationQuery.Page, paginationQuery.Limit)
	}

	return paginationQuery
}

type PaginationResponse struct {
	Page       uint64   // Current page number.
	TotalPages uint64   // Total number of pages.
	PagesLeft  uint64   // Number of pages remaining.
	TotalItems uint64   // Total number of items.
	ItemsLeft  uint64   // Number of items remaining on the current page.
	Limit      uint64   // Maximum items per page.
	OrderBy    string   // Field used for ordering.
	SortOrder  string   // Sorting direction ("asc" for ascending, "desc" for descending).
	PageLinks  []string // Array of page links.
	BaseURL    string   // Base URL for pagination.
}

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

func calculateTotalPages(totalItems, limit uint64) uint64 {
	totalPages := float64(totalItems) / float64(limit)
	return uint64(math.Ceil(totalPages))
}

func calculateItemsLeft(page, totalItems, limit uint64) uint64 {
	return totalItems - (page * limit)
}

// generatePageLinks generates the page links for the pagination response.
func generatePageLinks(paginationResponse PaginationResponse, amountOfPages uint64) []string {
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
func buildPageLink(paginationResponse PaginationResponse, pageNumber uint64) string {
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
	baseURL := paginationResponse.BaseURL
	return fmt.Sprintf("%s?%s", baseURL, queryParams)
}
