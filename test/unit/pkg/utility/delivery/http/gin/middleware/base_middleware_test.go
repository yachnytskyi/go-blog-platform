package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

func setupSecureHeadersMiddlewareConfig() *config.ApplicationConfig {
	config := mock.NewMockConfig()
	config.Security.ContentSecurityPolicyHeader.Key = test.ContentSecurityPolicyHeader
	config.Security.ContentSecurityPolicyHeader.Value = test.ContentSecurityPolicyValue
	config.Security.StrictTransportSecurityHeader.Key = test.StrictTransportSecurityHeader
	config.Security.StrictTransportSecurityHeader.Value = test.StrictTransportSecurityValue
	config.Security.XContentTypeOptionsHeader.Key = test.XContentTypeOptionsHeader
	config.Security.XContentTypeOptionsHeader.Value = test.XContentTypeOptionsValue
	return config
}

func setupCSPMiddlewareConfig() *config.ApplicationConfig {
	config := mock.NewMockConfig()
	config.Security.ContentSecurityPolicyHeaderFull.Key = test.ContentSecurityPolicyFullHeader
	config.Security.ContentSecurityPolicyHeaderFull.Value = test.ContentSecurityPolicyFullValue
	return config
}

func setupValidateInputMiddlewareConfig() *config.ApplicationConfig {
	config := mock.NewMockConfig()
	config.Security.AllowedHTTPMethods = []string{http.MethodGet, http.MethodPost}
	config.Security.AllowedContentTypes = []string{test.ContentTypeJSON}
	return config
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
	router := gin.Default()
	config := setupSecureHeadersMiddlewareConfig()
	router.Use(middleware.SecureHeadersMiddleware(config))

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
	router := gin.Default()
	config := setupCSPMiddlewareConfig()
	router.Use(middleware.CSPMiddleware(config))

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

func TestRateLimitMiddlewareWithDifferentLimits(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	config := mock.NewMockConfig()

	config.Security.RateLimit = 1
	router.Use(middleware.RateLimitMiddleware(config))

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
	router := gin.Default()
	config := mock.NewMockConfig()
	config.Security.RateLimit = 3
	router.Use(middleware.RateLimitMiddleware(config))

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

func TestValidateInputMiddlewareAcceptsValidRequest(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	config := setupValidateInputMiddlewareConfig()
	logger := mock.NewMockLogger()
	router.Use(middleware.ValidateInputMiddleware(config, logger))

	router.POST(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, test.TestURL, nil)
	request.Header.Set(constants.ContentType, test.ContentTypeJSON)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
}

func TestValidateInputMiddlewareRejectsInvalidMethod(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	config := setupValidateInputMiddlewareConfig()
	logger := mock.NewMockLogger()
	router.Use(middleware.ValidateInputMiddleware(config, logger))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, test.TestURL, nil)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, test.EqualMessage)
}

func TestValidateInputMiddlewareRejectsInvalidContentType(t *testing.T) {
	t.Parallel()
	router := gin.Default()
	config := setupValidateInputMiddlewareConfig()
	logger := mock.NewMockLogger()
	router.Use(middleware.ValidateInputMiddleware(config, logger))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, test.TestURL, nil)
	request.Header.Set(constants.ContentType, test.InvalidContentType)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, test.EqualMessage)
}

// func TestTimeoutMiddlewareCompletesWithinTimeout(t *testing.T) {
// 	t.Parallel()
// 	router := gin.Default()
// 	logger := mock.NewMockLogger()
// 	router.Use(middleware.TimeoutMiddleware(logger))

// 	router.GET(test.TestURL, func(ginContext *gin.Context) {
// 		time.Sleep(10 * time.Millisecond)
// 		ginContext.String(http.StatusOK, constants.Success)
// 	})

// 	recorder := httptest.NewRecorder()
// 	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
// 	router.ServeHTTP(recorder, request)

// 	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
// 	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
// }

// func TestTimeoutMiddlewareTimesOut(t *testing.T) {
// 	t.Parallel()
// 	router := gin.Default()
// 	logger := mock.NewMockLogger()
// 	router.Use(middleware.TimeoutMiddleware(logger))

// 	router.GET(test.TestURL, func(ginContext *gin.Context) {
// 		time.Sleep(constants.DefaultContextTimer + 100*time.Millisecond)
// 		ginContext.String(http.StatusOK, constants.Success)
// 	})

// 	recorder := httptest.NewRecorder()
// 	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
// 	router.ServeHTTP(recorder, request)

// 	assert.Equal(t, http.StatusBadGateway, recorder.Code, test.EqualMessage)
// 	assert.Contains(t, recorder.Body.String(), constants.InternalErrorNotification, test.EqualMessage)
// }

func TestLoggerMiddlewareLogsIncomingAndOutgoingRequests(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware(mockLogger))

	router.GET(test.TestURL, func(c *gin.Context) {
		c.String(http.StatusOK, constants.Success)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)

	assert.NotNil(t, mockLogger.LastInfo, test.DataNotNilMessage)
}
