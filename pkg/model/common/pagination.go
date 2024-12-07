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

func NewPaginationQuery(page, limit, orderBy, sortOrder, baseURL string) PaginationQuery {
	intPage := getPage(page)
	intLimit := getLimit(limit)

	return PaginationQuery{
		Page:      intPage,
		Limit:     intLimit,
		OrderBy:   getOrderBy(orderBy),
		SortOrder: getSortOrder(sortOrder),
		Skip:      getSkip(intPage, intLimit),
		BaseURL:   baseURL,
	}
}

func getPage(page string) int {
	intPage, stringConversionError := strconv.ParseInt(page, 0, 0)
	if validator.IsError(stringConversionError) || intPage < 1 {
		return constants.DefaultPageInteger
	}

	return int(intPage)
}

func getLimit(limit string) int {
	intLimit, stringConversionError := strconv.ParseInt(limit, 0, 0)
	if validator.IsError(stringConversionError) || intLimit < 1 {
		return constants.DefaultLimitInteger
	}
	if intLimit > constants.MaxItemsPerPage {
		return constants.DefaultLimitInteger
	}

	return int(intLimit)
}

func getOrderBy(orderBy string) string {
	if orderBy == "" {
		return constants.DefaultOrderBy
	}

	return orderBy
}

func getSortOrder(sortOrder string) string {
	if sortOrder != constants.SortAscend && sortOrder != constants.SortDescend {
		return constants.DefaultSortOrder
	}

	return sortOrder
}

func getSkip(page, limit int) int {
	return (page - 1) * limit
}

func SetCorrectPage(paginationQuery PaginationQuery) PaginationQuery {
	if paginationQuery.TotalItems <= paginationQuery.Skip {
		paginationQuery.Page = calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
		paginationQuery.Skip = getSkip(paginationQuery.Page, paginationQuery.Limit)
	}

	return paginationQuery
}

type PaginationResponse struct {
	Page       int      // Current page number.
	TotalPages int      // Total number of pages.
	PagesLeft  int      // Number of pages remaining.
	TotalItems int      // Total number of items.
	ItemsLeft  int      // Number of items remaining on the current page.
	PageStart  int      // The index of the first item displayed on the current page.
	PageEnd    int      // The index of the last item displayed on the current page.
	Limit      int      // Maximum items per page.
	OrderBy    string   // Field used for ordering.
	SortOrder  string   // Sorting direction ("asc" for ascending, "desc" for descending).
	PageLinks  []string // Array of page links.
}

func NewPaginationResponse(paginationQuery PaginationQuery) PaginationResponse {
	totalPages := calculateTotalPages(paginationQuery.TotalItems, paginationQuery.Limit)
	paginationResponse := PaginationResponse{
		Page:       paginationQuery.Page,
		TotalPages: totalPages,
		PagesLeft:  totalPages - paginationQuery.Page,
		TotalItems: paginationQuery.TotalItems,
		ItemsLeft:  calculateItemsLeft(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit),
		PageStart:  getPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit),
		PageEnd:    getPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit),
		Limit:      paginationQuery.Limit,
		OrderBy:    paginationQuery.OrderBy,
		SortOrder:  paginationQuery.SortOrder,
	}

	paginationResponse.PageLinks = generatePageLinks(paginationResponse, paginationQuery.BaseURL)
	return paginationResponse
}

func calculateTotalPages(totalItems, limit int) int {
	if totalItems == 0 {
		return 1
	}

	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

func calculateItemsLeft(page, totalItems, limit int) int {
	if totalItems <= page*limit {
		return 0
	}

	return totalItems - (page * limit)
}

// Calculates the index of the first item displayed on the given page.
func getPageStart(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}
	if page > 1 {
		return (page-1)*limit + 1
	}

	return page
}

// Calculates the index of the last item displayed on the given page.
func getPageEnd(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}

	end := page * limit
	if end > totalItems {
		return totalItems
	}

	return end
}

// generatePageLinks generates the page links for the pagination response.
func generatePageLinks(paginationResponse PaginationResponse, baseURL string) []string {
	// Preallocate memory for the pageLinks slice based on amountOfPageLinks, adding space for potential first and last page links.
	pageLinks := make([]string, 0, constants.DefaultAmountOfPageLinks+2)

	// Calculate the start page to display, centering around the current page.
	startPage := paginationResponse.Page - (constants.DefaultAmountOfPageLinks / 2)
	if startPage < 1 {
		startPage = 1
	}

	// Calculate the end page based on the adjusted start page and amountOfPageLinks.
	endPage := startPage + constants.DefaultAmountOfPageLinks
	if startPage == 1 {
		endPage++
	}

	// Ensure the endPage does not exceed the total number of pages, and adjust startPage accordingly.
	if endPage >= paginationResponse.TotalPages {
		endPage = paginationResponse.TotalPages
		startPage = endPage - constants.DefaultAmountOfPageLinks - 1
		if startPage < 1 {
			startPage = 1
		}
	}

	// If the first page is not included, add it to the pageLinks slice.
	if startPage > 1 {
		pageLinks = append(pageLinks, buildPageLink(paginationResponse, 1, baseURL))
	}

	// Append links for pages within the calculated range, excluding the current page.
	for index := startPage; index <= endPage; index++ {
		pageLinks = append(pageLinks, buildPageLink(paginationResponse, index, baseURL))
	}

	// If the last page is not included, add it to the pageLinks slice.
	if endPage < paginationResponse.TotalPages {
		pageLinks = append(pageLinks, buildPageLink(paginationResponse, paginationResponse.TotalPages, baseURL))
	}

	return pageLinks
}

// buildPageLink builds the page link with given page number.
func buildPageLink(paginationResponse PaginationResponse, pageNumber int, baseURL string) string {
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
