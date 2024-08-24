package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

func setupSecureHeadersMiddlewareConfig() mock.MockConfig {
	config := mock.NewMockConfig()
	configInstance := config.GetConfig()

	configInstance.Security.ContentSecurityPolicyHeader.Key = test.ContentSecurityPolicyHeader
	configInstance.Security.ContentSecurityPolicyHeader.Value = test.ContentSecurityPolicyValue
	configInstance.Security.StrictTransportSecurityHeader.Key = test.StrictTransportSecurityHeader
	configInstance.Security.StrictTransportSecurityHeader.Value = test.StrictTransportSecurityValue
	configInstance.Security.XContentTypeOptionsHeader.Key = test.XContentTypeOptionsHeader
	configInstance.Security.XContentTypeOptionsHeader.Value = test.XContentTypeOptionsValue
	return config
}

func TestRequestIDMiddlewareAddsCorrelationIDWhenAbsent(t *testing.T) {
	router := gin.Default()
	recorder := httptest.NewRecorder()
	router.Use(middleware.RequestIDMiddleware())

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)
	correlationID := recorder.Header().Get(constants.CorrelationIDHeader)

	assert.NotEmpty(t, correlationID, test.NotFailureMessage)
}

func TestRequestIDMiddlewareUsesExistingCorrelationID(t *testing.T) {
	router := gin.Default()
	recorder := httptest.NewRecorder()
	router.Use(middleware.RequestIDMiddleware())

	existingCorrelationID := "existing-correlation-id"
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.CorrelationIDHeader, existingCorrelationID)
	router.ServeHTTP(recorder, request)

	correlationID := recorder.Header().Get(constants.CorrelationIDHeader)
	assert.Equal(t, existingCorrelationID, correlationID, test.FailureMessage)
}

func TestRequestIDMiddlewareSetsCorrelationIDInContext(t *testing.T) {
	router := gin.Default()
	recorder := httptest.NewRecorder()

	router.GET(test.TestURL, func(c *gin.Context) {
		correlationID := c.GetString(constants.CorrelationIDHeader)
		c.String(http.StatusOK, correlationID)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)

	correlationID := recorder.Header().Get(constants.CorrelationIDHeader)
	assert.Equal(t, correlationID, recorder.Body.String(), test.FailureMessage)
}

func TestSecureHeadersMiddlewareSetsHeadersCorrectly(t *testing.T) {
	router := gin.Default()
	config := setupSecureHeadersMiddlewareConfig()
	router.Use(middleware.SecureHeadersMiddleware(config.ApplicationConfig))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	router.ServeHTTP(recorder, request)

	headers := recorder.Header()
	assert.Equal(t, test.ContentSecurityPolicyValue, headers.Get(test.ContentSecurityPolicyHeader), test.FailureMessage)
	assert.Equal(t, test.StrictTransportSecurityValue, headers.Get(test.StrictTransportSecurityHeader), test.FailureMessage)
	assert.Equal(t, test.XContentTypeOptionsValue, headers.Get(test.XContentTypeOptionsHeader), test.FailureMessage)
}
