package common

import (
	"strconv"

	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	DefaultPage     = "1"
	DefaultLimit    = "10"
	maxItemsPerPage = 100
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
		Skip:    GetSkip(page, limit),
	}
}

func GetPage(page string) int {
	intPage, stringConversionError := strconv.Atoi(page)
	if validator.IsErrorNotNil(stringConversionError) {
		intPage, _ = strconv.Atoi(DefaultPage)
	}
	if isPaginationValueNotValid(intPage) {
		intPage, _ = strconv.Atoi(DefaultPage)
	}
	return intPage
}

func GetLimit(limit string) int {
	intLimit, stringConversionError := strconv.Atoi(limit)
	if validator.IsErrorNotNil(stringConversionError) {
		intLimit, _ = strconv.Atoi(DefaultLimit)
	}
	if isPaginationValueNotValid(intLimit) {
		intLimit, _ = strconv.Atoi(DefaultLimit)
	}
	return intLimit
}

func GetOrderBy(orderBy string) string {
	return orderBy
}

func GetSkip(page, limit int) int {
	return (page - 1) * limit
}

func isPaginationValueNotValid(data int) bool {
	if data == 0 || data < 0 || data > maxItemsPerPage {
		return true
	}
	return false
}
