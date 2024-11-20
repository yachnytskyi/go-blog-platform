package common

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	"github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	defaultPage       = 1
	defaultLimit      = 10
	defaultTotalPages = 10
	orderBy           = "name"
	totalItems        = 100
	invalidSortOrder  = "invalid"
	baseURL           = "http://localhost:8080/api/users"
)

func calculateSkip(page, limit int) int {
	return (page - 1) * limit
}

func TestNewPaginationQueryFirstPage(t *testing.T) {
	t.Parallel()
	expectedSkip := calculateSkip(defaultPage, defaultLimit)
	expectedTotalItems := 0
	result := common.NewPaginationQuery(strconv.Itoa(int(defaultPage)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryLastPage(t *testing.T) {
	t.Parallel()
	page := 20
	expectedSkip := calculateSkip(page, defaultLimit)
	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryValidInputs(t *testing.T) {
	t.Parallel()
	page := 2
	expectedSkip := calculateSkip(page, defaultLimit)
	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryValidLimitOne(t *testing.T) {
	t.Parallel()
	page := 2
	limit := "1"
	limitInt := 1
	expectedSkip := calculateSkip(page, limitInt)
	expectedTotalItems := 0
	result := common.NewPaginationQuery(strconv.Itoa(int(page)), limit, orderBy, constants.SortDescend, baseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, limitInt, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortDescend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryZeroValues(t *testing.T) {
	t.Parallel()
	page := "0"
	limit := "0"
	pageInt := 0
	limitInt := 0
	expectedSkip := calculateSkip(pageInt, limitInt)
	expectedTotalItems := 0
	result := common.NewPaginationQuery(page, limit, orderBy, constants.SortAscend, baseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryNegativeValues(t *testing.T) {
	t.Parallel()
	page := -5
	limit := -10
	expectedSkip := calculateSkip(defaultPage, defaultLimit)
	expectedTotalItems := 0
	result := common.NewPaginationQuery(strconv.Itoa(int(page)), strconv.Itoa(int(limit)), orderBy, constants.SortAscend, baseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryEmptyData(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery("", "", "", "", baseURL)
	expectedSkip := calculateSkip(defaultPage, defaultLimit)
	expectedTotalItems := 0

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryInvalidSortOrder(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery(constants.DefaultPage, constants.DefaultLimit, orderBy, invalidSortOrder, baseURL)
	expectedSkip := calculateSkip(defaultPage, defaultLimit)
	expectedTotalItems := 0

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestSetCorrectPageNoCorrectionNeeded(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery(constants.DefaultPage, constants.DefaultLimit, "", "", baseURL)
	paginationQuery.TotalItems = 0
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, paginationQuery, result, test.EqualMessage)
}

func TestSetCorrectPageLimitIsMoreThanTotalItems(t *testing.T) {
	t.Parallel()
	totalItems := 5
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(defaultPage), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, paginationQuery, result, test.EqualMessage)

}

func TestSetCorrectPageSkipIsLessThanTotalItems(t *testing.T) {
	t.Parallel()
	page := 1
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	paginationQuery.TotalItems = totalItems
	paginationQuery.Skip = calculateSkip(page, paginationQuery.Limit)

	expectedPage := "1"
	expected := common.NewPaginationQuery(expectedPage, constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = calculateSkip(expected.Page, expected.Limit)
	expected.TotalItems = totalItems
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestSetCorrectPageSkipIsMoreThanTotalItems(t *testing.T) {
	t.Parallel()
	page := 15
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	paginationQuery.Skip = calculateSkip(page, paginationQuery.Limit)
	paginationQuery.TotalItems = totalItems

	expectedPage := "10"
	expected := common.NewPaginationQuery(expectedPage, constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = calculateSkip(expected.Page, expected.Limit)
	expected.TotalItems = totalItems
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestNewPaginationResponseEmptyData(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery("", "", "", "", "")
	paginationQuery.TotalItems = 0
	expectedAmount := 0
	expectedPageLinks := []string{
		"?page=1&limit=10&order_by=created_at&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultPage, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedAmount, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedAmount, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedAmount, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, paginationQuery.Limit, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseSinglePage(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery(constants.DefaultPage, constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems / 10
	expectedTotalItems := totalItems / 10
	expectedAmount := 0
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultPage, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedAmount, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedAmount, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, paginationQuery.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseMiddlePage(t *testing.T) {
	t.Parallel()
	page := 5
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems
	expectedTotalPages := 10
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=2&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=3&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=4&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=5&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=6&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=7&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, totalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, totalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseFirstPage(t *testing.T) {
	t.Parallel()
	page := 1
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=2&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=3&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=4&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=5&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=6&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=7&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseSeventhPage(t *testing.T) {
	t.Parallel()
	page := 7
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortDescend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=2&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=3&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=4&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=5&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=6&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=7&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=descend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortDescend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseEighthPage(t *testing.T) {
	t.Parallel()
	page := 8
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=3&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=4&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=5&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=6&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=7&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=13&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseTwelvethPage(t *testing.T) {
	t.Parallel()
	page := 12
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=7&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=13&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=14&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=15&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=16&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=17&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationFourteenthPage(t *testing.T) {
	t.Parallel()
	page := 13
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortDescend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=8&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=13&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=14&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=15&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=16&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=17&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=18&limit=10&order_by=name&sort_order=descend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=descend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortDescend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationEighteenthPage(t *testing.T) {
	t.Parallel()
	page := 18
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=13&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=14&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=15&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=16&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=17&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=18&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=19&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationLastPage(t *testing.T) {
	t.Parallel()
	page := 20
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageLinks := []string{
		"http://localhost:8080/api/users?page=1&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=9&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=10&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=11&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=12&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=13&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=14&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=15&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=16&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=17&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=18&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=19&limit=10&order_by=name&sort_order=ascend",
		"http://localhost:8080/api/users?page=20&limit=10&order_by=name&sort_order=ascend",
	}
	result := common.NewPaginationResponse(paginationQuery)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}
