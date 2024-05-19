// Package middleware provides a set of Gin middleware functions for enhancing the security and functionality of your web application.
//
// SecureHeadersMiddleware: Adds secure headers to HTTP responses.
// CSPMiddleware: Adds Content Security Policy (CSP) headers to HTTP responses.
// RateLimitMiddleware: Implements rate limiting to control the number of requests from clients.
// LoggingMiddleware: Logs incoming requests and outgoing responses with additional context.
// ValidateInput: Validates input based on HTTP methods and content types.
package middleware

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/common"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// SecureHeadersMiddleware adds secure headers to HTTP responses.
// Returns a Gin middleware handler function.
func SecureHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve application configuration for security headers.
		securityConfig := config.AppConfig.Security

		// Set secure headers.
		c.Header(securityConfig.ContentSecurityPolicyHeader.Key, securityConfig.ContentSecurityPolicyHeader.Value)
		c.Header(securityConfig.StrictTransportSecurityHeader.Key, securityConfig.StrictTransportSecurityHeader.Value)
		c.Header(securityConfig.XContentTypeOptionsHeader.Key, securityConfig.XContentTypeOptionsHeader.Value)

		// Continue to the next middleware or handler in the chain.
		c.Next()
	}
}

// CSPMiddleware adds Content Security Policy (CSP) headers to HTTP responses.
// Returns a Gin middleware handler function.
func CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve application configuration for CSP headers.
		securityConfig := config.AppConfig.Security

		// Set the Content Security Policy header.
		c.Writer.Header().Set(
			securityConfig.ContentSecurityPolicyHeaderFull.Key,
			securityConfig.ContentSecurityPolicyHeaderFull.Value,
		)

		// Continue to the next middleware or handler in the chain.
		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting to control the number of requests from clients.
// Returns a Gin middleware handler function.
func RateLimitMiddleware() gin.HandlerFunc {
	// Retrieve application configuration for rate limiting.
	securityConfig := config.AppConfig.Security

	// Create an instance of limiter.ExpirableOptions.
	limiterOptions := &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}

	// Create a rate limiter with the specified options.
	limiter := tollbooth.NewLimiter(securityConfig.RateLimit, limiterOptions)
	return func(c *gin.Context) {
		tollbooth_gin.LimitHandler(limiter)(c)
		c.Next()
	}
}

// LoggingMiddleware logs incoming requests and outgoing responses with additional context.
// Returns a Gin middleware handler function.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Log information about the incoming request.
		incomingRequest := commonModel.HTTPLog{
			Location:      location + "LoggingMiddleware",
			Time:          start,
			RequestMethod: c.Request.Method,
			RequestURL:    c.Request.URL.Path,
			ClientIP:      c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
		}

		// Continue processing the request.
		c.Next()

		// Log information about the outgoing response.
		outgoingResponse := commonModel.HTTPLog{
			Location:       location + "LoggingMiddleware",
			Time:           start,
			RequestMethod:  c.Request.Method,
			RequestURL:     c.Request.URL.Path,
			ClientIP:       c.ClientIP(),
			UserAgent:      c.Request.UserAgent(),
			ResponseStatus: c.Writer.Status(),
			Duration:       time.Since(start),
		}

		logging.Logger(incomingRequest)
		logging.Logger(outgoingResponse)
	}
}

// ValidateInput allows specific HTTP methods and checks for the "application/json" or "application/grpc" content type.
// Returns a Gin middleware handler function.
func ValidateInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the content type from the request header.
		contentType := c.GetHeader("Content-Type")

		// Retrieve application configuration for allowed requests.
		securityConfig := config.AppConfig.Security

		// Check if the request method is in the list of allowed methods.
		if validator.IsSliceNotContains(securityConfig.AllowedHTTPMethods, c.Request.Method) {
			// Reject requests with disallowed HTTP methods.
			httpRequestError := httpError.NewHttpRequestErrorView(contentType, constants.InvalidHTTPMethodNotification)
			abortWithStatusJSON(c, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Check if the content type is in the list of allowed content types.
		if validator.IsStringNotEmpty(contentType) && validator.IsSliceNotContains(securityConfig.AllowedContentTypes, contentType) {
			// Reject requests with disallowed content types.
			httpRequestError := httpError.NewHttpRequestErrorView(contentType, constants.InvalidContentTypeNotification)
			abortWithStatusJSON(c, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Continue processing the request.
		c.Next()
	}
}
