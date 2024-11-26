package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const ()

// RequestIDMiddleware adds a request ID to requests and responses.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		requestID := ginContext.GetHeader(constants.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ginContext.Set(constants.RequestIDHeader, requestID)
		ginContext.Writer.Header().Set(constants.RequestIDHeader, requestID)
		ginContext.Next()
	}
}

// SecureHeadersMiddleware adds secure headers to HTTP responses.
func SecureHeadersMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginContext.Header(config.Security.ContentSecurityPolicyHeader.Key, config.Security.ContentSecurityPolicyHeader.Value)
		ginContext.Header(config.Security.StrictTransportSecurityHeader.Key, config.Security.StrictTransportSecurityHeader.Value)
		ginContext.Header(config.Security.XContentTypeOptionsHeader.Key, config.Security.XContentTypeOptionsHeader.Value)
		ginContext.Next()
	}
}

// CSPMiddleware adds Content Security Policy (CSP) headers to HTTP responses.
func CSPMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginContext.Writer.Header().Set(
			config.Security.ContentSecurityPolicyHeaderFull.Key,
			config.Security.ContentSecurityPolicyHeaderFull.Value,
		)

		ginContext.Next()
	}
}

// RateLimitMiddleware implements rate limiting to control the number of requests from clients.
func RateLimitMiddleware(config *config.ApplicationConfig) gin.HandlerFunc {
	limiterOptions := &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}

	limiter := tollbooth.NewLimiter(config.Security.RateLimit, limiterOptions)
	return func(ginContext *gin.Context) {
		tollbooth_gin.LimitHandler(limiter)(ginContext)
		ginContext.Next()
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
			httpRequestError := httpError.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedHTTPMethods", ginContext.Request.Method, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, http.StatusBadRequest)
			return
		}

		// Check if the content type is in the list of allowed content types.
		if contentType != "" && validator.IsSliceNotContains(config.Security.AllowedContentTypes, contentType) {
			allowedContentTypes := strings.Join(config.Security.AllowedContentTypes, ", ")
			notification := constants.InvalidHTTPMethodNotification + allowedContentTypes
			httpRequestError := httpError.NewHTTPRequestError(location+"ValidateInputMiddleware.AllowedContentTypes", contentType, notification)
			abortWithStatusJSON(ginContext, logger, httpRequestError, http.StatusBadRequest)
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
			httpInternalError := httpError.NewHTTPInternalError(location+"TimeOutMiddleware", ctx.Err().Error())
			abortWithStatusJSON(ginContext, logger, httpInternalError, http.StatusBadGateway)
		}
	}
}

// LoggerMiddleware logs incoming requests and outgoing responses with additional context.
func LoggerMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		start := time.Now()
		requestID := ginContext.GetString(constants.RequestIDHeader)

		httpIncomingLog := delivery.NewHTTPIncomingLog(
			location+"LoggerMiddleware",
			requestID,
			ginContext.Request.Method,
			ginContext.Request.URL.Path,
			ginContext.ClientIP(),
			ginContext.Request.UserAgent(),
		)
		ginContext.Next()

		httpOutgoingLog := delivery.NewHTTPOutgoingLog(
			location+"LoggerMiddleware",
			requestID,
			ginContext.Request.Method,
			ginContext.Request.URL.Path,
			ginContext.ClientIP(),
			ginContext.Request.UserAgent(),
			ginContext.Writer.Status(),
			time.Since(start),
		)

		logger.Info(httpIncomingLog)
		logger.Info(httpOutgoingLog)
	}
}
