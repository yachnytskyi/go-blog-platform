package common

import (
	"fmt"

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
	page := ginContext.DefaultQuery(constants.Page, constants.DefaultPage)
	limit := ginContext.DefaultQuery(constants.Limit, constants.DefaultLimit)
	orderBy := ginContext.DefaultQuery(constants.OrderBy, constants.DefaultOrderBy)
	sortOrder := ginContext.DefaultQuery(constants.SortOrder, constants.DefaultSortOrder)

	// Convert and validate the extracted values.
	convertedPage := commonModel.GetPage(page)
	convertedLimit := commonModel.GetLimit(limit)

	// Extract the scheme from the request.
	scheme := constants.HTTP
	if ginContext.Request.TLS != nil {
		scheme = constants.HTTPS
	}

	// Extract the full URL path including query parameters.
	rawURL := ginContext.Request.URL

	// Get the scheme and host from the original request.
	host := ginContext.Request.Host

	// Construct the full URL by combining scheme, host, and the URL path.
	baseURL := fmt.Sprintf("%s://%s%s", scheme, host, rawURL.Path)

	// Create and return a commonModel.PaginationQuery struct.
	return commonModel.NewPaginationQuery(convertedPage, convertedLimit, orderBy, sortOrder, baseURL)
}
