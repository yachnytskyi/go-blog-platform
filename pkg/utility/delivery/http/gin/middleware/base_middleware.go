// Package middleware provides a set of Gin middleware functions for enhancing the security and functionality of your web application.
//
// SecureHeadersMiddleware: Adds secure headers to HTTP responses.
// CSPMiddleware: Adds Content Security Policy (CSP) headers to HTTP responses.
// RateLimitMiddleware: Implements rate limiting to control the number of requests from clients.
// LoggingMiddleware: Logs incoming requests and outgoing responses with additional context.
// ValidateInput: Validates input based on HTTP methods and content types.
package middleware

import (
	"fmt"
	"strings"
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
		securityConfig := config.GetSecurityConfig()

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
		securityConfig := config.GetSecurityConfig()

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
	securityConfig := config.GetSecurityConfig()

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
		httpIncomingLog := commonModel.NewHTTPIncomingLog(
			location+"LoggingMiddleware",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		// Continue processing the request.
		c.Next()

		// Log information about the outgoing response.
		httpOutgoingLog := commonModel.NewHTTPOutgoingLog(
			location+"LoggingMiddleware",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			time.Since(start),
		)

		logging.Logger(httpIncomingLog)
		logging.Logger(httpOutgoingLog)
	}
}

// ValidateInputMiddleware allows specific HTTP methods and checks for the content type.
// Returns a Gin middleware handler function.
func ValidateInputMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Get the content type from the request header.
		contentType := ginContext.GetHeader(constants.ContentType)

		// Retrieve application configuration for allowed requests.
		securityConfig := config.GetSecurityConfig()

		// Check if the request method is in the list of allowed methods.
		if validator.IsSliceNotContains(securityConfig.AllowedHTTPMethods, ginContext.Request.Method) {
			allowedMethods := strings.Join(config.GetSecurityConfig().AllowedHTTPMethods, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedMethods

			// Reject requests with disallowed HTTP methods.
			httpRequestError := httpError.NewHttpRequestErrorView(location+"ValidateInputMiddleware.AllowedHTTPMethods", ginContext.Request.Method, notification)
			abortWithStatusJSON(ginContext, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Check if the content type is in the list of allowed content types.
		if validator.IsStringNotEmpty(contentType) && validator.IsSliceNotContains(securityConfig.AllowedContentTypes, contentType) {
			allowedContentTypes := strings.Join(config.GetSecurityConfig().AllowedContentTypes, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedContentTypes

			// Reject requests with disallowed content types.
			httpRequestError := httpError.NewHttpRequestErrorView(location+"ValidateInputMiddleware.AllowedContentTypes", contentType, notification)
			abortWithStatusJSON(ginContext, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Continue processing the request.
		ginContext.Next()
	}
}

// EnforceHTTPMethod ensures that the incoming request uses the specified HTTP method.
// If the request method does not match the specified method, it responds with a 405 Method Not Allowed status.
// Returns a Gin middleware handler function.
func EnforceHTTPMethod(method string) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Check if the request method matches the specified method.
		if ginContext.Request.Method != method {
			// Create a notification message indicating the allowed method.
			notification := fmt.Sprintf("Method %s not allowed. You can only use %s method", ginContext.Request.Method, method)

			// Create an HTTP request error view with the location, method, and notification.
			httpRequestError := httpError.NewHttpRequestErrorView(location+"ValidateInputMiddleware.EnforceHTTPMethod", ginContext.Request.Method, notification)

			// Abort the request and respond with the error and status code 405 (Method Not Allowed).
			abortWithStatusJSON(ginContext, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Continue to the next middleware or handler in the chain if the method matches.
		ginContext.Next()
	}
}
