package common

import (
	"math"
	"strconv"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// PaginationQuery represents the parameters for paginating a list of items.
type PaginationQuery struct {
	Page      int    // Page number to retrieve.
	Limit     int    // Maximum number of items per page.
	OrderBy   string // Field to order by.
	SortOrder string // Sorting direction ("asc" for ascending, "desc" for descending).
	Skip      int    // Number of items to skip for pagination.
}

// NewPaginationQuery creates a new PaginationQuery with the given parameters.
func NewPaginationQuery(page, limit int, orderBy, sortOrder string) PaginationQuery {
	return PaginationQuery{
		Page:      page,
		Limit:     limit,
		OrderBy:   orderBy,
		SortOrder: sortOrder,
		Skip:      getSkip(page, limit),
	}
}

// GetPage converts a string to an integer representing the page number.
func GetPage(page string) int {
	intPage, stringConversionError := strconv.Atoi(page)
	if validator.IsErrorNotNil(stringConversionError) {
		intPage, _ = strconv.Atoi(constants.DefaultPage)
	}
	if validator.IsIntegerZeroOrNegative(intPage) {
		intPage, _ = strconv.Atoi(constants.DefaultPage)
	}
	return intPage
}

// GetLimit converts a string to an integer representing the maximum items per page.
func GetLimit(limit string) int {
	intLimit, stringConversionError := strconv.Atoi(limit)
	if validator.IsErrorNotNil(stringConversionError) {
		intLimit, _ = strconv.Atoi(constants.DefaultLimit)
	}
	if isLimitNotValid(intLimit) {
		intLimit, _ = strconv.Atoi(constants.DefaultLimit)
	}
	return intLimit
}

// getSkip calculates the number of items to skip based on the current page and limit.
func getSkip(page, limit int) int {
	if validator.IsIntegerZero(page) {
		return page
	}
	return (page - 1) * limit
}

// isLimitNotValid checks if a limit value is valid.
func isLimitNotValid(data int) bool {
	if data == 0 || data < 0 || data > constants.MaxItemsPerPage {
		return true
	}
	return false
}

// SetCorrectPage adjusts the PaginationQuery to ensure it's valid, especially when there are not enough items to reach the current page.
func SetCorrectPage(totalItems int, paginationQuery PaginationQuery) PaginationQuery {
	if totalItems <= paginationQuery.Skip {
		paginationQuery.Page = getTotalPages(totalItems, paginationQuery.Limit)
		paginationQuery.Skip = getSkip(paginationQuery.Page, paginationQuery.Limit)
	}
	return paginationQuery
}

// PaginationResponse represents information about the current page, total pages, and more for a paginated list.
type PaginationResponse struct {
	Page       int    // Current page number.
	TotalPages int    // Total number of pages.
	PagesLeft  int    // Number of pages remaining.
	TotalItems int    // Total number of items.
	ItemsLeft  int    // Number of items remaining on the current page.
	Limit      int    // Maximum items per page.
	OrderBy    string // Field used for ordering.
}

// NewPaginationResponse creates a new PaginationResponse with the given parameters.
func NewPaginationResponse(page, totalItems, limit int, orderBy string) PaginationResponse {
	return PaginationResponse{
		Page:       page,
		TotalPages: getTotalPages(totalItems, limit),
		PagesLeft:  getPagesLeft(page, totalItems, limit),
		TotalItems: totalItems,
		ItemsLeft:  getItemsLeft(page, totalItems, limit),
		Limit:      limit,
		OrderBy:    orderBy,
	}
}

// getTotalPages calculates the total number of pages based on total items and limit.
func getTotalPages(totalItems, limit int) int {
	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

// getPagesLeft calculates the number of pages remaining.
func getPagesLeft(page, totalItems, limit int) int {
	return getTotalPages(totalItems, limit) - page
}

// getItemsLeft calculates the number of items remaining on the current page.
func getItemsLeft(page, totalItems, limit int) int {
	if validator.IsIntegerNegative(totalItems - (page * limit)) {
		return 0
	}
	return totalItems - (page * limit)
}
