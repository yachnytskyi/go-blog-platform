package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

// ParsePaginationQuery parses and extracts pagination parameters from the provided Gin context.
func ParsePaginationQuery(ginContext *gin.Context) commonModel.PaginationQuery {
	page := ginContext.DefaultQuery(constants.Page, constants.DefaultPage)
	limit := ginContext.DefaultQuery(constants.Limit, constants.DefaultLimit)
	orderBy := ginContext.DefaultQuery(constants.OrderBy, constants.DefaultOrderBy)
	sortOrder := ginContext.DefaultQuery(constants.SortOrder, constants.DefaultSortOrder)

	mappedPage := commonModel.GetPage(page)
	mappedLimit := commonModel.GetLimit(limit)

	scheme := constants.HTTP
	if ginContext.Request.TLS != nil {
		scheme = constants.HTTPS
	}

	rawURL := ginContext.Request.URL
	host := ginContext.Request.Host
	baseURL := fmt.Sprintf("%s://%s%s", scheme, host, rawURL.Path)

	return commonModel.NewPaginationQuery(
		mappedPage,
		mappedLimit,
		orderBy,
		sortOrder,
		baseURL,
	)
}
