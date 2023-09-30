package common

import (
	"math"
	"strconv"

	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	zero = 0
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
		Skip:    getSkip(page, limit),
	}
}

func GetPage(page string) int {
	intPage, stringConversionError := strconv.Atoi(page)
	if validator.IsErrorNotNil(stringConversionError) {
		intPage, _ = strconv.Atoi(constant.DefaultPage)
	}
	if validator.IsIntegerZeroOrNegative(intPage) {
		intPage, _ = strconv.Atoi(constant.DefaultPage)
	}
	return intPage
}

func GetLimit(limit string) int {
	intLimit, stringConversionError := strconv.Atoi(limit)
	if validator.IsErrorNotNil(stringConversionError) {
		intLimit, _ = strconv.Atoi(constant.DefaultLimit)
	}
	if isLimitNotValid(intLimit) {
		intLimit, _ = strconv.Atoi(constant.DefaultLimit)
	}
	return intLimit
}

func getSkip(page, limit int) int {
	return (page - 1) * limit
}

func isLimitNotValid(data int) bool {
	if data == zero || data < zero || data > constant.MaxItemsPerPage {
		return true
	}
	return false
}

func SetCorrectPage(totalItems int, paginationQuery PaginationQuery) PaginationQuery {
	if totalItems <= paginationQuery.Skip {
		paginationQuery.Page = getTotalPages(totalItems, paginationQuery.Limit)
		paginationQuery.Skip = getSkip(paginationQuery.Page, paginationQuery.Limit)
	}
	return paginationQuery
}

type PaginationResponse struct {
	Page       int
	TotalPages int
	PagesLeft  int
	TotalItems int
	ItemsLeft  int
	Limit      int
	OrderBy    string
}

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

func getTotalPages(totalItems, limit int) int {
	totalPages := float64(totalItems) / float64(limit)
	return int(math.Ceil(totalPages))
}

func getPagesLeft(page, totalItems, limit int) int {
	return getTotalPages(totalItems, limit) - page
}

func getItemsLeft(page, totalItems, limit int) int {
	if validator.IsIntegerNegative(totalItems - (page * limit)) {
		return zero
	}
	return totalItems - (page * limit)
}
