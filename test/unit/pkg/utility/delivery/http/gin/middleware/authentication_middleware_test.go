package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
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

func setupAuthenticationMiddlewareConfig() mock.MockConfig {
	config := mock.NewMockConfig()
	configInstance := config.GetConfig()
	configInstance.AccessToken.PublicKey = test.PublicKey
	configInstance.AccessToken.PrivateKey = test.PrivateKey
	configInstance.AccessToken.ExpiredIn = constants.PasswordResetTokenExpirationTime
	return config
}

func getValidToken(location string) common.Result[string] {
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	validToken := utility.GenerateJWTToken(
		logger,
		location,
		config.ApplicationConfig.AccessToken.PrivateKey,
		config.ApplicationConfig.AccessToken.ExpiredIn,
		tokenPayload,
	)
	return validToken
}

func getExpiredTokenForAuthenticationMiddleware(location string) common.Result[string] {
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	expiredToken := utility.GenerateJWTToken(
		logger,
		location+"TestAuthenticationMiddlewareExpiredToken",
		config.ApplicationConfig.AccessToken.PrivateKey,
		-config.ApplicationConfig.AccessToken.ExpiredIn,
		tokenPayload,
	)
	return expiredToken
}

func TestAuthenticationMiddlewareCookieValidToken(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareCookieValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Nil(t, validToken.Error, test.NotFailureMessage)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, constants.Success, recorder.Body.String())
}

func TestAuthenticationMiddlewareCookieInvalidTokenValue(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: "invalid token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareCookieEmptyTokenValue(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: ""})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareCookieMultipleTokens(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareCookieMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: validToken.Data + ", " + validToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareCookieExpiredToken(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForAuthenticationMiddleware(location + "TestAuthenticationMiddlewareCookieExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: expiredToken.Data})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareHeaderValidToken(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareHeaderValidToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Nil(t, validToken.Error, test.NotFailureMessage)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, constants.Success, recorder.Body.String())
}

func TestAuthenticationMiddlewareHeaderInvalidTokenValue(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+"invalid token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareHeaderEmptyTokenValue(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareHeaderMultipleTokens(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	validToken := getValidToken(location + "TestAuthenticationMiddlewareHeaderMultipleTokens")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+validToken.Data+", "+validToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareHeaderExpiredToken(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	expiredToken := getExpiredTokenForAuthenticationMiddleware(location + "TestAuthenticationMiddlewareHeaderExpiredToken")
	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+expiredToken.Data)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}

func TestAuthenticationMiddlewareNoToken(t *testing.T) {
	router := gin.Default()
	logger := mock.NewMockLogger()
	config := setupAuthenticationMiddlewareConfig()

	router.GET(test.TestURL, middleware.AuthenticationMiddleware(config, logger), func(ginContext *gin.Context) {
		ginContext.String(http.StatusOK, constants.Success)
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.JSONEq(t, test.NotLoggedInMessage, recorder.Body.String())
}
