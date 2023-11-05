package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

// ParsePaginationQuery parses and extracts pagination parameters from the provided Gin context.
// It takes the Gin context as input and returns a commonModel.PaginationQuery struct.
//
// Parameters:
// - ginContext: The Gin context containing the HTTP request.
//
// Returns:
// - commonModel.PaginationQuery: The parsed pagination parameters.
func ParsePaginationQuery(ginContext *gin.Context) commonModel.PaginationQuery {
	// Extract the page, limit, orderBy, and sortOrder query parameters from the Gin context.
	page := ginContext.DefaultQuery("page", constants.DefaultPage)
	limit := ginContext.DefaultQuery("limit", constants.DefaultLimit)
	orderBy := ginContext.DefaultQuery("order-by", constants.DefaultOrderBy)
	sortOrder := ginContext.DefaultQuery("sort-order", constants.DefaultSortOrder)

	// Convert and validate the extracted values.
	convertedPage := commonModel.GetPage(page)
	convertedLimit := commonModel.GetLimit(limit)

	// Create and return a commonModel.PaginationQuery struct.
	return commonModel.NewPaginationQuery(convertedPage, convertedLimit, orderBy, sortOrder)
}
