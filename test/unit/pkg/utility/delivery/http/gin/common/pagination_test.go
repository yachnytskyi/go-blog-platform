package common

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

func TestParsePaginationQueryDefaultValues(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		common.GetPage(constants.DefaultPage),
		common.GetLimit(constants.DefaultLimit),
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQueryCustomValues(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=2&limit=10&order_by=name&sort_rder=descend"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		common.GetPage("2"),
		common.GetLimit("10"),
		"name",
		"descend",
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQueryEmptyParameters(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=&limit="
	result := delivery.ParsePaginationQuery(ginContext)
	page, _ := strconv.ParseUint(constants.DefaultPage, 0, 0)
	limit, _ := strconv.ParseUint(constants.DefaultLimit, 0, 0)
	expected := common.NewPaginationQuery(
		page,
		limit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQueryHTTPScheme(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
		TLS:  &tls.ConnectionState{},
	}

	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		common.GetPage(constants.DefaultPage),
		common.GetLimit(constants.DefaultLimit),
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		"https://localhost:8080/test",
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQuerySpecialCharacters(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=1&limit=10&order_by=special!@#$%^&*()&sort_order=asc"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		common.GetPage("1"),
		common.GetLimit("10"),
		"created_at",
		"descend",
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQueryInvalidParameters(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "pageInvalid=-48912481&limitInvalid=-9124829148&orderBy=invalidField&sortOrder=asc"
	page, _ := strconv.ParseUint(constants.DefaultPage, 0, 0)
	limit, _ := strconv.ParseUint(constants.DefaultLimit, 0, 0)
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		page,
		limit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}

func TestParsePaginationQueryExtraParameters(t *testing.T) {
	t.Parallel()

	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=3&limit=15&order_by=date&sort_order=ascend&unexpectedParam=value"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		common.GetPage("3"),
		common.GetLimit("15"),
		"date",
		"ascend",
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result)
}
