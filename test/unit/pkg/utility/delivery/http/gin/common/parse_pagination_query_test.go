package common

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	testUrl = "https://localhost:8080/test"
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
		constants.DefaultPage,
		constants.DefaultLimit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestParsePaginationQueryCustomValues(t *testing.T) {
	t.Parallel()
	page := "2"
	orderBy := "name"
	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=2&limit=10&order_by=name&sort_order=descend"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		page,
		constants.DefaultLimit,
		orderBy,
		constants.SortDescend,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestParsePaginationQueryNegativeValues(t *testing.T) {
	t.Parallel()
	orderBy := "name"
	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=-42187412&limit=-214124&order_by=name&sort_order=ascend"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		constants.DefaultPage,
		constants.DefaultLimit,
		orderBy,
		constants.SortAscend,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
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
	expected := common.NewPaginationQuery(
		constants.DefaultPage,
		constants.DefaultLimit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
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
		constants.DefaultPage,
		constants.DefaultLimit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		testUrl,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestParsePaginationQueryInvalidParameters(t *testing.T) {
	t.Parallel()
	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "pageInvalid=-48912481&limitInvalid=-9124829148&orderBy=invalidField&sortOrder=asc"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		constants.DefaultPage,
		constants.DefaultLimit,
		constants.DefaultOrderBy,
		constants.DefaultSortOrder,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
}

func TestParsePaginationQueryExtraParameters(t *testing.T) {
	t.Parallel()
	page := "3"
	limit := "15"
	orderBy := "date"
	ginContext, _ := gin.CreateTestContext(nil)
	ginContext.Request = &http.Request{
		Host: test.Localhost,
		URL:  &url.URL{Path: test.TestURL},
	}

	ginContext.Request.URL.RawQuery = "page=3&limit=15&order_by=date&sort_order=ascend&unexpectedParam=value"
	result := delivery.ParsePaginationQuery(ginContext)
	expected := common.NewPaginationQuery(
		page,
		limit,
		orderBy,
		constants.SortAscend,
		test.LocalhostTest,
	)

	assert.Equal(t, expected, result, test.EqualMessage)
}
