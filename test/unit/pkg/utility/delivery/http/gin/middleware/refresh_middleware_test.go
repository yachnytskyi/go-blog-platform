package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

func setupRefreshTokenMiddlewareConfig() *config.ApplicationConfig {
	mockConfig := mock.NewMockConfig()
	mockConfig.RefreshToken.PublicKey = test.PublicKey
	mockConfig.RefreshToken.PrivateKey = test.PrivateKey
	mockConfig.RefreshToken.ExpiredIn = constants.PasswordResetTokenExpirationTime
	return mockConfig
}

func getExpiredTokenForRefreshTokenMiddleware(location string) common.Result[string] {
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	expiredToken := utility.GenerateJWTToken(
		mockLogger,
		location+"TestAuthenticationMiddlewareExpiredToken",
		mockConfig.RefreshToken.PrivateKey,
		-mockConfig.RefreshToken.ExpiredIn,
		tokenPayload,
	)

	return expiredToken
}

func TestRefreshTokenMiddlewareCookieValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestRefreshTokenMiddlewareCookieValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.RefreshTokenValue, Value: validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.NoError(t, validToken.Error, test.ErrorNilMessage)
	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestRefreshTokenMiddlewareHeaderValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.NoError(t, validToken.Error, test.ErrorNilMessage)
	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareCookieInvalidTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.RefreshTokenValue, Value: "invalid token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareCookieEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.RefreshTokenValue, Value: ""})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareCookieMultipleTokens(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestRefreshTokenMiddlewareMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.RefreshTokenValue, Value: validToken.Data + ", " + validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareCookieExpiredToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForRefreshTokenMiddleware(location + "TestRefreshTokenMiddlewareExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.RefreshTokenValue, Value: expiredToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderInvalidTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+"invalid token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderMultipleTokens(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestRefreshTokenMiddlewareHeaderMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data+", "+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderExpiredToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForRefreshTokenMiddleware(location + "TestRefreshTokenMiddlewareHeaderExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+expiredToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestRefreshTokenMiddlewareHeaderNoToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupRefreshTokenMiddlewareConfig()
	router := gin.Default()
	router.Use(middleware.RefreshTokenMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}
