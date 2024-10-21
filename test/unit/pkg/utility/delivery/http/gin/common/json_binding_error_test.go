package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
	ginMock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/delivery/http/gin"
)

const (
	location       = "test.unit.pkg.utility.delivery.gin.common."
	invalidRequest = "invalid request"
)

func setupTestRouter(location string, mockLogger *mock.MockLogger) *gin.Engine {
	router := ginMock.NewMockGinEngine()
	router.POST(test.TestURL, func(ginContext *gin.Context) {
		var json map[string]interface{}
		shouldBindJSONError := ginContext.ShouldBindJSON(&json)
		if validator.IsError(shouldBindJSONError) {
			common.HandleJSONBindingError(ginContext, mockLogger, location, shouldBindJSONError)
			return
		}
	})

	return router
}

func makeTestRequest(router *gin.Engine) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, test.TestURL, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestHandleJSONBindingError(t *testing.T) {
	mockLogger := mock.NewMockLogger()
	router := setupTestRouter(location+"TestHandleJSONBindingError", mockLogger)
	recorder := makeTestRequest(router)

	expectedResponse := `{
	"error": {
		"notification": "` + constants.InternalErrorNotification + `"
	},
	"status": "` + constants.Fail + `"
}` + "\n"

	expectedLocation := location + "TestHandleJSONBindingError.ShouldBindJSON"
	expectedErrorMessage := fmt.Sprintf(test.ExpectedErrorMessageFormat, expectedLocation, invalidRequest)

	assert.IsType(t, httpError.HTTPInternalError{}, mockLogger.LastError, test.EqualMessage)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, test.EqualMessage)
	assert.JSONEq(t, expectedResponse, recorder.Body.String(), test.EqualMessage)
	assert.Equal(t, expectedErrorMessage, mockLogger.LastError.Error(), test.EqualMessage)
}
