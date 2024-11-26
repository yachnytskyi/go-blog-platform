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
	defaultPage       = 1
	defaultLimit      = 10
	defaultTotalPages = 10
	orderBy           = "name"
	totalItems        = 100
	invalidSortOrder  = "invalid"
	baseURL           = "http://localhost:8080/api/users"
)

// getSkipTestHelper is a test-specific copy of the getSkipTest logic.
func getSkipTestHelper(page, limit int) int {
	return (page - 1) * limit
}

// getPageStartTestHelper is a test-specific copy of the getPageStart logic.
func getPageStartTestHelper(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}

	if page > 1 {
		return (page-1)*limit + 1
	}

	return page
}

// getPageEndTestHelper is a test-specific copy of the getPageEnd logic.
func getPageEndTestHelper(page, totalItems, limit int) int {
	if totalItems == 0 {
		return 0
	}

	end := page * limit
	if end > totalItems {
		return totalItems
	}

	return end
}

// generatePageLinksTestHelper is a test-specific copy of the generatePageLinks logic.
func generatePageLinksTestHelper(paginationResponse common.PaginationResponse, baseURL string) []string {
	pageLinks := make([]string, 0, constants.DefaultAmountOfPageLinks+2)
	startPage := paginationResponse.Page - (constants.DefaultAmountOfPageLinks / 2)
	if startPage < 1 {
		startPage = 1
	}

	endPage := startPage + constants.DefaultAmountOfPageLinks
	if startPage == 1 {
		endPage++
	}

	if endPage >= paginationResponse.TotalPages {
		endPage = paginationResponse.TotalPages
		startPage = endPage - constants.DefaultAmountOfPageLinks - 1
		if startPage < 1 {
			startPage = 1
		}
	}

	if startPage > 1 {
		pageLinks = append(pageLinks, buildPageLinkTestHelper(paginationResponse, 1, baseURL))
	}

	for index := startPage; index <= endPage; index++ {
		pageLinks = append(pageLinks, buildPageLinkTestHelper(paginationResponse, index, baseURL))
	}

	if endPage < paginationResponse.TotalPages {
		pageLinks = append(pageLinks, buildPageLinkTestHelper(paginationResponse, paginationResponse.TotalPages, baseURL))
	}

	return pageLinks
}

// buildPageLinkTestHelper is a test-specific copy of the buildPageLink logic.
func buildPageLinkTestHelper(paginationResponse common.PaginationResponse, pageNumber int, baseURL string) string {
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

	return fmt.Sprintf("%s?%s", baseURL, queryParams)
}

func TestNewPaginationQueryFirstPage(t *testing.T) {
	t.Parallel()
	result := common.NewPaginationQuery(strconv.Itoa(int(defaultPage)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getSkipTestHelper(defaultPage, defaultLimit)
	expectedTotalItems := 0

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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getSkipTestHelper(page, defaultLimit)

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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	expectedSkip := getSkipTestHelper(page, defaultLimit)

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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), limit, orderBy, constants.SortDescend, baseURL)
	expectedSkip := getSkipTestHelper(page, limitInt)
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
	expectedSkip := getSkipTestHelper(pageInt, limitInt)
	expectedTotalItems := 0

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

	result := common.NewPaginationQuery(strconv.Itoa(int(page)), strconv.Itoa(int(limit)), orderBy, constants.SortAscend, baseURL)
	expectedSkip := getSkipTestHelper(defaultPage, defaultLimit)
	expectedTotalItems := 0

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
	expectedSkip := getSkipTestHelper(defaultPage, defaultLimit)
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
	expectedSkip := getSkipTestHelper(defaultPage, defaultLimit)
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
	paginationQuery.Skip = getSkipTestHelper(page, paginationQuery.Limit)

	result := common.SetCorrectPage(paginationQuery)
	expectedPage := "1"
	expected := common.NewPaginationQuery(expectedPage, constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = getSkipTestHelper(expected.Page, expected.Limit)
	expected.TotalItems = totalItems

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestSetCorrectPageSkipIsMoreThanTotalItems(t *testing.T) {
	t.Parallel()
	page := 15
	paginationQuery := common.NewPaginationQuery(strconv.Itoa(page), constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	paginationQuery.Skip = getSkipTestHelper(page, paginationQuery.Limit)
	paginationQuery.TotalItems = totalItems

	result := common.SetCorrectPage(paginationQuery)
	expectedPage := "10"
	expected := common.NewPaginationQuery(expectedPage, constants.DefaultLimit, constants.DefaultOrderBy, constants.DefaultSortOrder, baseURL)
	expected.Skip = getSkipTestHelper(expected.Page, expected.Limit)
	expected.TotalItems = totalItems

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestNewPaginationResponseEmptyData(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery("", "", "", "", "")
	paginationQuery.TotalItems = 0

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalItems := 0
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultPage, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
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
	expectedTotalPages := 10
	expectedTotalItems := paginationQuery.TotalItems
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, constants.DefaultOrderBy, result.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.DefaultSortOrder, result.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)

}

func TestNewPaginationResponseSinglePage(t *testing.T) {
	t.Parallel()
	paginationQuery := common.NewPaginationQuery(constants.DefaultPage, constants.DefaultLimit, orderBy, constants.SortAscend, baseURL)
	paginationQuery.TotalItems = totalItems / 10

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalItems := totalItems / 10
	expectedItemsPagesLeft := 0
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, defaultPage, result.Page, test.EqualMessage)
	assert.Equal(t, defaultPage, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedItemsPagesLeft, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedItemsPagesLeft, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 10
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, totalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, totalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
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

	result := common.NewPaginationResponse(paginationQuery)
	expectedTotalPages := 20
	expectedTotalItems := totalItems * 2
	expectedPageStart := getPageStartTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageEnd := getPageEndTestHelper(paginationQuery.Page, paginationQuery.TotalItems, paginationQuery.Limit)
	expectedPageLinks := generatePageLinksTestHelper(result, paginationQuery.BaseURL)

	assert.Equal(t, page, result.Page, test.EqualMessage)
	assert.Equal(t, expectedTotalPages, result.TotalPages, test.EqualMessage)
	assert.Equal(t, expectedTotalPages-paginationQuery.Page, result.PagesLeft, test.EqualMessage)
	assert.Equal(t, expectedTotalItems, result.TotalItems, test.EqualMessage)
	assert.Equal(t, expectedTotalItems-paginationQuery.Page*paginationQuery.Limit, result.ItemsLeft, test.EqualMessage)
	assert.Equal(t, expectedPageStart, result.PageStart, test.EqualMessage)
	assert.Equal(t, expectedPageEnd, result.PageEnd, test.EqualMessage)
	assert.Equal(t, defaultLimit, result.Limit, test.EqualMessage)
	assert.Equal(t, orderBy, paginationQuery.OrderBy, test.EqualMessage)
	assert.Equal(t, constants.SortAscend, paginationQuery.SortOrder, test.EqualMessage)
	assert.ElementsMatch(t, expectedPageLinks, result.PageLinks, test.EqualMessage)
}
