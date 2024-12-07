package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// ParsePaginationQuery parses and extracts pagination parameters from the provided Gin context.
func ParsePaginationQuery(ginContext *gin.Context) common.PaginationQuery {
	page := ginContext.DefaultQuery(constants.Page, constants.DefaultPage)
	limit := ginContext.DefaultQuery(constants.Limit, constants.DefaultLimit)
	orderBy := ginContext.DefaultQuery(constants.OrderBy, constants.DefaultOrderBy)
	sortOrder := ginContext.DefaultQuery(constants.SortOrder, constants.DefaultSortOrder)
	scheme := constants.HTTP
	if validator.IsValueNotEmpty(ginContext.Request.TLS) {
		scheme = constants.HTTPS
	}

	return common.NewPaginationQuery(
		page,
		limit,
		orderBy,
		sortOrder,
		fmt.Sprintf("%s://%s%s", scheme, ginContext.Request.Host, ginContext.Request.URL.Path))
}
