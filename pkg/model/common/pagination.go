package common

import (
	"math"
	"strconv"

	"github.com/yachnytskyi/golang-mongo-grpc/config"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type PaginationQuery struct {
	Page    int
	Limit   int
	OrderBy string
	Skip    int
}

func NewPaginationQuery(page, limit int, orderBy string) PaginationQuery {
	return PaginationQuery{
		Page:    page,
		Limit:   limit,
		OrderBy: orderBy,
		Skip:    (page - 1) * limit,
	}
}

func GetPage(page string) int {
	intPage, stringConversionError := strconv.Atoi(page)
	if validator.IsErrorNotNil(stringConversionError) {
		intPage, _ = strconv.Atoi(config.DefaultPage)
	}
	if validator.IsIntegerZeroOrLess(intPage) {
		intPage, _ = strconv.Atoi(config.DefaultPage)
	}
	return intPage
}

func GetLimit(limit string) int {
	intLimit, stringConversionError := strconv.Atoi(limit)
	if validator.IsErrorNotNil(stringConversionError) {
		intLimit, _ = strconv.Atoi(config.DefaultLimit)
	}
	if isLimitNotValid(intLimit) {
		intLimit, _ = strconv.Atoi(config.DefaultLimit)
	}
	return intLimit
}

func GetOrderBy(orderBy string) string {
	return orderBy
}

func isLimitNotValid(data int) bool {
	if data == 0 || data < 0 || data > config.MaxItemsPerPage {
		return true
	}
	return false
}

type PaginationResponse struct {
	CurrentPage int
	TotalPages  int
	PagesLeft   int
	TotalItems  int
	Limit       int
	OrderBy     string
}

func NewPaginationResponse(currentPage, totalItems, limit int, orderBy string) PaginationResponse {
	return PaginationResponse{
		CurrentPage: currentPage,
		TotalPages:  GetTotalPages(totalItems, limit),
		PagesLeft:   GetPagesLeft(totalItems, limit, currentPage),
		TotalItems:  totalItems,
		Limit:       limit,
		OrderBy:     orderBy,
	}
}

func GetTotalPages(totalItems, limit int) int {
	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

func GetPagesLeft(totalItems, limit, currentPage int) int {
	return GetTotalPages(totalItems, limit) - currentPage
}
