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
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

const (
	location = "test.unit.pkg.utility.delivery.http.gin.middleware."
)

var (
	tokenPayload = domain.UserTokenPayload{
		UserID: "12345",
		Role:   "user",
	}
)

func setupAuthenticationMiddlewareConfig() *config.ApplicationConfig {
	mockConfig := mock.NewMockConfig()
	mockConfig.AccessToken.PublicKey = test.PublicKey
	mockConfig.AccessToken.PrivateKey = test.PrivateKey
	mockConfig.AccessToken.ExpiredIn = constants.PasswordResetTokenExpirationTime
	return mockConfig
}

func getValidToken(location string) common.Result[string] {
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	validToken := utility.GenerateJWTToken(
		mockLogger,
		location,
		mockConfig.AccessToken.PrivateKey,
		mockConfig.AccessToken.ExpiredIn,
		tokenPayload,
	)

	return validToken
}

func getExpiredTokenForAuthenticationMiddleware(location string) common.Result[string] {
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	expiredToken := utility.GenerateJWTToken(
		mockLogger,
		location+"TestAuthenticationMiddlewareExpiredToken",
		mockConfig.AccessToken.PrivateKey,
		-mockConfig.AccessToken.ExpiredIn,
		tokenPayload,
	)

	return expiredToken
}

func TestAuthenticationMiddlewareCookieValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareCookieValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.NoError(t, validToken.Error, test.NotFailureMessage)
	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareCookieInvalidTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: "invalid token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareCookieEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: ""})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareCookieMultipleTokens(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareCookieMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: validToken.Data + ", " + validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareCookieExpiredToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForAuthenticationMiddleware(location + "TestAuthenticationMiddlewareCookieExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: expiredToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareHeaderValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareHeaderValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.NoError(t, validToken.Error, test.NotFailureMessage)
	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, constants.Success, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareHeaderInvalidTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
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

func TestAuthenticationMiddlewareHeaderEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
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

func TestAuthenticationMiddlewareHeaderMultipleTokens(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareHeaderMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data+", "+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareHeaderExpiredToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForAuthenticationMiddleware(location + "TestAuthenticationMiddlewareHeaderExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+expiredToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, delivery.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAuthenticationMiddlewareNoToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	mockConfig := setupAuthenticationMiddlewareConfig()

	router := gin.Default()
	router.Use(middleware.AuthenticationMiddleware(mockConfig, mockLogger))
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
