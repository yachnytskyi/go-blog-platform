package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/common"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const ()

// RequestIDMiddleware adds a correlation ID to requests and responses.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(constants.CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		c.Set(constants.CorrelationIDHeader, correlationID)
		c.Writer.Header().Set(constants.CorrelationIDHeader, correlationID)
		c.Next()
	}
}

// SecureHeadersMiddleware adds secure headers to HTTP responses.
func SecureHeadersMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(config.Security.ContentSecurityPolicyHeader.Key, config.Security.ContentSecurityPolicyHeader.Value)
		c.Header(config.Security.StrictTransportSecurityHeader.Key, config.Security.StrictTransportSecurityHeader.Value)
		c.Header(config.Security.XContentTypeOptionsHeader.Key, config.Security.XContentTypeOptionsHeader.Value)
		c.Next()
	}
}

// CSPMiddleware adds Content Security Policy (CSP) headers to HTTP responses.
func CSPMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(
			config.Security.ContentSecurityPolicyHeaderFull.Key,
			config.Security.ContentSecurityPolicyHeaderFull.Value,
		)

		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting to control the number of requests from clients.
func RateLimitMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	limiterOptions := &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}

	limiter := tollbooth.NewLimiter(config.Security.RateLimit, limiterOptions)
	return func(c *gin.Context) {
		tollbooth_gin.LimitHandler(limiter)(c)
		c.Next()
	}
}

// ValidateInputMiddleware allows specific HTTP methods and checks for the content type.
func ValidateInputMiddleware(config *config.ApplicationConfig, logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		contentType := ginContext.GetHeader(constants.ContentType)

		// Check if the request method is in the list of allowed methods.
		if validator.IsSliceNotContains(config.Security.AllowedHTTPMethods, ginContext.Request.Method) {
			allowedMethods := strings.Join(config.Security.AllowedHTTPMethods, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedMethods
			httpRequestError := http.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedHTTPMethods", ginContext.Request.Method, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, constants.StatusBadRequest)
			return
		}

		// Check if the content type is in the list of allowed content types.
		if contentType != "" && validator.IsSliceNotContains(config.Security.AllowedContentTypes, contentType) {
			allowedContentTypes := strings.Join(config.Security.AllowedContentTypes, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedContentTypes
			httpRequestError := http.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedContentTypes", contentType, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, constants.StatusBadRequest)
			return
		}

		ginContext.Next()
	}
}

// TimeoutMiddleware sets a timeout for each request.
func TimeoutMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		ginContext.Request = ginContext.Request.WithContext(ctx)
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			ginContext.Next()
		}()

		select {
		case <-ch:
		case <-ctx.Done():
			httpInternalError := http.NewHTTPInternalError(location+"TimeOutMiddleware", ctx.Err().Error())
			abortWithStatusJSON(ginContext, logger, httpInternalError, constants.StatusBadGateway)
		}
	}
}

// LoggerMiddleware logs incoming requests and outgoing responses with additional context.
func LoggerMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		correlationID := c.GetString(constants.CorrelationIDHeader)

		httpIncomingLog := common.NewHTTPIncomingLog(
			location+"LoggerMiddleware",
			correlationID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		c.Next()

		httpOutgoingLog := common.NewHTTPOutgoingLog(
			location+"LoggerMiddleware",
			correlationID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			time.Since(start),
		)

		logger.Info(httpIncomingLog)
		logger.Info(httpOutgoingLog)
	}
}
