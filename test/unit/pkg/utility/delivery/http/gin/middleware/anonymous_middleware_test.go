package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

func TestAnonymousMiddlewareNoToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, test.SuccessResponse, recorder.Body.String(), test.EqualMessage)
}

func TestAnonymousMiddlewareHeaderEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, test.EqualMessage)
	assert.Equal(t, test.SuccessResponse, recorder.Body.String(), test.EqualMessage)
}

func TestAnonymousMiddlewareCookieEmptyTokenValue(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: ""})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, httpError.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusForbidden, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAnonymousMiddlewareCookieValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.AddCookie(&http.Cookie{Name: constants.AccessTokenValue, Value: "valid token"})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, httpError.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusForbidden, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAnonymousMiddlewareHeaderValidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+"valid token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, httpError.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusForbidden, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}

func TestAnonymousMiddlewareHeaderInvalidToken(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	router := gin.Default()
	router.Use(middleware.AnonymousMiddleware(mockLogger))
	router.GET(test.TestURL, func(ginContext *gin.Context) {
		ginContext.JSON(http.StatusOK, gin.H{test.Message: constants.Success})
	})

	request := httptest.NewRequest(http.MethodGet, test.TestURL, nil)
	request.Header.Set(constants.Authorization, constants.Bearer+" ")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.IsType(t, httpError.HTTPAuthorizationError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusForbidden, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, test.AlreadyLoggedInMessage, recorder.Body.String(), test.EqualMessage)
}
