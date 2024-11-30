package common

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	"github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	orderBy          = "name"
	totalItems       = 100
	totalPages       = 10
	baseURL          = "http://localhost:8080/api/users"
	invalidSortOrder = "invalid"
)

func getExpectedSkip(page, limit int) int {
	return (page - 1) * limit
}

func getExpectedPageStart(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}

	if page > 1 {
		return (page-1)*limit + 1
	}

	return page
}

func getExpectedPageEnd(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}

	end := page * limit
	if end > totalItems {
		return totalItems
	}

	return end
}

// generateExpectedPageLinks creates paginated URLs based on the given pagination response and page numbers.
func generateExpectedPageLinks(paginationResponse common.PaginationResponse, pages []int) []string {
	var pageLinks []string
	for _, page := range pages {
		link := fmt.Sprintf(
			"%s?page=%d&limit=%d&order_by=%s&sort_order=%s",
			baseURL,
			page,
			paginationResponse.Limit,
			paginationResponse.OrderBy,
			paginationResponse.SortOrder,
		)
		pageLinks = append(pageLinks, link)
	}

	return pageLinks
}

func TestNewPaginationQueryFirstPage(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery(constants.DefaultPage, constants.Limit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getExpectedSkip(constants.DefaultPageInteger, constants.DefaultLimitInteger)
	expectedTotalItems := 0

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryLastPage(t *testing.T) {
	t.Parallel()
	page := 20

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.Limit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getExpectedSkip(page, constants.DefaultLimitInteger)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryValidInputs(t *testing.T) {
	t.Parallel()
	page := 2

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.Limit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getExpectedSkip(page, constants.DefaultLimitInteger)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), limit, orderBy, constants.SortDescend, baseURL)
	expectedSkip := getExpectedSkip(page, limitInt)
	expectedTotalItems := 0

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

	result := common.NewPaginationQuery(page, limit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getExpectedSkip(pageInt, limitInt)
	expectedTotalItems := 0

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), strconv.Itoa(int(limit)), orderBy, constants.SortAscend, baseURL)
	expectedSkip := getExpectedSkip(constants.DefaultPageInteger, constants.DefaultLimitInteger)
	expectedTotalItems := 0

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryEmptyData(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery("", "", "", "", baseURL)
	expectedSkip := getExpectedSkip(constants.DefaultPageInteger, constants.DefaultLimitInteger)
	expectedTotalItems := 0

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestNewPaginationQueryInvalidSortOrder(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery(constants.Page, constants.Limit, orderBy, invalidSortOrder, baseURL)
	expectedSkip := getExpectedSkip(constants.DefaultPageInteger, constants.DefaultLimitInteger)
	expectedTotalItems := 0

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.Equal(t, expectedSkip, result.Skip, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, baseURL, result.BaseURL, test.EqualMessage)
}

func TestSetCorrectPageNoCorrectionNeeded(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery(constants.Page, constants.Limit, "", "", baseURL)
	paginationQuery.TotalItems = 0
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, paginationQuery, result, test.EqualMessage)
}

func TestSetCorrectPageLimitIsMoreThanTotalItems(t *testing.T) {
	t.Parallel()
	totalItems := 5
	paginationQuery := common.NewPaginationQuery(constants.DefaultPage, constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems
	result := common.SetCorrectPage(paginationQuery)

	assert.Equal(t, paginationQuery, result, test.EqualMessage)

}

func TestSetCorrectPageSkipIsLessThanTotalItems(t *testing.T) {
	t.Parallel()
	page := 1
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	paginationQuery.TotalItems = totalItems
	paginationQuery.Skip = getExpectedSkip(page, paginationQuery.Limit)

	result := common.SetCorrectPage(paginationQuery)
	expectedPage := "1"
	expected := common.NewPaginationQuery(expectedPage, constants.Limit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = getExpectedSkip(expected.Page, expected.Limit)
	expected.TotalItems = totalItems

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestSetCorrectPageSkipIsMoreThanTotalItems(t *testing.T) {
	t.Parallel()
	page := 15
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	paginationQuery.Skip = getExpectedSkip(page, paginationQuery.Limit)
	paginationQuery.TotalItems = totalItems

	result := common.SetCorrectPage(paginationQuery)
	expectedPage := "10"
	expected := common.NewPaginationQuery(expectedPage, constants.Limit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = getExpectedSkip(expected.Page, expected.Limit)
	expected.TotalItems = totalItems

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestNewPaginationResponseEmptyData(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery("", "", "", "", "")
	paginationQuery.TotalItems = 0
	paginationQuery.BaseURL = baseURL

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalItems := 0
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultPageInteger, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponsePageEnd(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery("", "", "", "", baseURL)
	paginationQuery.TotalItems = 95

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalItems := paginationQuery.TotalItems
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, totalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, totalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseSinglePage(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery(constants.DefaultPage, constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = 1

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalItems := 1
	expectedItemsPagesLeft := 0
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, constants.DefaultPageInteger, result.Page, test.EqualMessage)
	assert.Equal(t, constants.DefaultPageInteger, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedItemsPagesLeft, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedItemsPagesLeft, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, paginationQuery.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseMiddlePage(t *testing.T) {
	t.Parallel()
	page := 5
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems

	result := common.NewPaginationResponse(paginationQuery)
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, totalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, totalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, totalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, totalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseFirstPage(t *testing.T) {
	t.Parallel()
	page := 1
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseSeventhPage(t *testing.T) {
	t.Parallel()
	page := 7
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortDescend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortDescend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseEighthPage(t *testing.T) {
	t.Parallel()
	page := 8
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationResponseTwelvethPage(t *testing.T) {
	t.Parallel()
	page := 12
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationFourteenthPage(t *testing.T) {
	t.Parallel()
	page := 13
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortDescend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortDescend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationEighteenthPage(t *testing.T) {
	t.Parallel()
	page := 18
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}

func TestNewPaginationLastPage(t *testing.T) {
	t.Parallel()
	page := 20
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.Limit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems * 2

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := totalPages * 2
	expectedTotalItems := totalItems * 2
	expectedPageStart := getExpectedPageStart(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getExpectedPageEnd(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPages := []int{1, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	expectedPageLinks := generateExpectedPageLinks(result, expectedPages)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultLimitInteger, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}
