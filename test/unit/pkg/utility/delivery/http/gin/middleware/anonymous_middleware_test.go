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

func TestAnonymousMiddlewareNoToken(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, test.SuccessResponse, recorder.Body.String())
}

func TestAnonymousMiddlewareCookieEmptyTokenValue(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: ""})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String())
}

func TestAnonymousMiddlewareCookieValidToken(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: "valid token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String())
}

func TestAnonymousMiddlewareHeaderEmptyTokenValue(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, test.SuccessResponse, recorder.Body.String())
}

func TestAnonymousMiddlewareHeaderValidToken(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+"valid token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String())
}

func TestAnonymousMiddlewareHeaderInvalidToken(t *testing.T) {
	router := gin.Default()
	mockLogger := mock.NewMockLogger()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+" ")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String())
}
