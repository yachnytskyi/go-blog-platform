package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

func setupSecureHeadersMiddlewareConfig() *config.ApplicationConfig {
	mockConfig := mock.NewMockConfig()
	mockConfig.Security.ContentSecurityPolicyHeader.Key = test.ContentSecurityPolicyHeader
	mockConfig.Security.ContentSecurityPolicyHeader.Value = test.ContentSecurityPolicyValue
	mockConfig.Security.StrictTransportSecurityHeader.Key = test.StrictTransportSecurityHeader
	mockConfig.Security.StrictTransportSecurityHeader.Value = test.StrictTransportSecurityValue
	mockConfig.Security.XContentTypeOptionsHeader.Key = test.XContentTypeOptionsHeader
	mockConfig.Security.XContentTypeOptionsHeader.Value = test.XContentTypeOptionsValue
	return mockConfig
}

func setupCSPMiddlewaremockConfig() *config.ApplicationConfig {
	mockConfig := mock.NewMockConfig()
	mockConfig.Security.ContentSecurityPolicyHeaderFull.Key = test.ContentSecurityPolicyFullHeader
	mockConfig.Security.ContentSecurityPolicyHeaderFull.Value = test.ContentSecurityPolicyFullValue
	return mockConfig
}

func setupValidateInputMiddlewareConfig() *config.ApplicationConfig {
	mockConfig := mock.NewMockConfig()
	mockConfig.Security.AllowedHTTPMethods = []string{http.MethodGet, http.MethodPost}
	mockConfig.Security.AllowedContentTypes = []string{test.ContentTypeJSON}
	return mockConfig
}

func TestRequestIDMiddlewareAddsRequestIDWhenAbsent(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	recorder := httptest.NewRecorder()
	router.Use(middleware.RequestIDMiddleware())

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	requestID := recorder.Header().Get(constants.RequestIDHeader)

	assert.NotEmpty(t, requestID, test.DataNotNilMessage)
}

func TestRequestIDMiddlewareUsesExistingRequestID(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	recorder := httptest.NewRecorder()
	router.Use(middleware.RequestIDMiddleware())

	existingRequestID := "existing-request-id"
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.RequestIDHeader, existingRequestID)
	router.ServeHTTP(recorder, request)
	requestID := recorder.Header().Get(constants.RequestIDHeader)

	assert.Equal(t, existingRequestID, requestID, test.EqualMessage)
}

func TestRequestIDMiddlewareSetsRequestIDInContext(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	recorder := httptest.NewRecorder()
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		requestID := ginContext.GetString(constants.RequestIDHeader)
		ginContext.String(http.StatusOK, requestID)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	requestID := recorder.Header().Get(constants.RequestIDHeader)

	assert.Equal(t, requestID, recorder.Body.String(), test.EqualMessage)
}

func TestSecureHeadersMiddlewareSetsHeadersCorrectly(t *testing.T) {
	t.Parallel()
	mockConfig := setupSecureHeadersMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.SecureHeadersMiddleware(mockConfig))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	headers := recorder.Header()

	assert.Equal(t, test.ContentSecurityPolicyValue, headers.Get(test.ContentSecurityPolicyHeader), test.EqualMessage)
	assert.Equal(t, test.StrictTransportSecurityValue, headers.Get(test.StrictTransportSecurityHeader), test.EqualMessage)
	assert.Equal(t, test.XContentTypeOptionsValue, headers.Get(test.XContentTypeOptionsHeader), test.EqualMessage)
}

func TestCSPMiddlewareSetsCSPHeader(t *testing.T) {
	t.Parallel()
	mockConfig := setupCSPMiddlewaremockConfig()
	router := gin.Default()
	router.Use(middleware.CSPMiddleware(mockConfig))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	header := recorder.Header().Get(test.ContentSecurityPolicyFullHeader)

	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
	assert.Equal(t, test.ContentSecurityPolicyFullValue, header, test.EqualMessage)
}

func TestValidateInputMiddlewareAcceptsValidRequest(t *testing.T) {
	t.Parallel()
	mockConfig := setupValidateInputMiddlewareConfig()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.ValidateInputMiddleware(mockConfig, mockLogger))
	router.POST(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, test.TestURL, nil)
	request.Header.Set(constants.ContentType, test.ContentTypeJSON)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
}

func TestLoggerMiddlewareLogsIncomingAndOutgoingRequests(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware(mockLogger))
	router.GET(test.TestURL, func(c *gin.Context) {
		c.String(http.StatusOK, constants.Success)
	})

	userAgentKey := "User-Agent"
	userAgentValue := "test-user-agent"
	expectedLocation := "pkg.utility.delivery.http.gin.middleware.LoggerMiddleware"
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(userAgentKey, userAgentValue)
	router.ServeHTTP(recorder, request)

	outgoingLog, _ := mockLogger.LastInfo.(delivery.HTTPOutgoingLog)
	assert.IsType(t, delivery.HTTPOutgoingLog{}, mockLogger.LastInfo, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
	assert.Equal(t, outgoingLog.Location, expectedLocation, test.EqualMessage)
	assert.NotEmpty(t, outgoingLog.RequestID, test.EqualMessage)
	assert.Equal(t, http.MethodGet, outgoingLog.RequestMethod, test.EqualMessage)
	assert.Equal(t, test.TestURL, outgoingLog.RequestURL, test.EqualMessage)
	assert.NotEmpty(t, outgoingLog.ClientIP, test.EqualMessage)
	assert.Equal(t, userAgentValue, outgoingLog.UserAgent, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
	assert.Equal(t, http.StatusOK, outgoingLog.ResponseStatus, test.EqualMessage)
	assert.NotZero(t, outgoingLog.Duration, test.EqualMessage)
}

func TestRateLimitMiddlewareWithDifferentLimits(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Security.RateLimit = 1
	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(mockConfig))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code, test.EqualMessage)
}

func TestRateLimitMiddlewareLimitsRequests(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Security.RateLimit = 3
	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(mockConfig))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)

	for i := 0; i < 4; i++ {
		recorder = httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
	}

	assert.Equal(t, http.StatusTooManyRequests, recorder.Code, test.EqualMessage)
}

func TestValidateInputMiddlewareRejectsInvalidMethod(t *testing.T) {
	t.Parallel()
	mockConfig := setupValidateInputMiddlewareConfig()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.ValidateInputMiddleware(mockConfig, mockLogger))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, test.TestURL, nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, test.EqualMessage)
}

func TestValidateInputMiddlewareRejectsInvalidContentType(t *testing.T) {
	t.Parallel()
	mockConfig := setupValidateInputMiddlewareConfig()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.ValidateInputMiddleware(mockConfig, mockLogger))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, test.TestURL, nil)
	request.Header.Set(constants.ContentType, test.InvalidContentType)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, test.EqualMessage)
}
