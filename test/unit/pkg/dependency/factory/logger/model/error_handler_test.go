package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/model"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	"github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location = "test.unit.pkg.dependency.factory.logger.model."

	field        = "Username"
	notification = "Some test notification"
	requestType = "GET"
)

func TestHandleErrorHTTPAuthorizationError(t *testing.T) {
	authorizationError := http.NewHTTPAuthorizationError(location+"TestHandleErrorHTTPAuthorizationError", notification)
	result := model.HandleError(authorizationError)

	jsonAuthorizationError, ok := result.(model.JSONAuthorizationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, authorizationError.Location, jsonAuthorizationError.Location, test.EqualMessage)
	assert.Equal(t, authorizationError.Notification, jsonAuthorizationError.Notification, test.EqualMessage)
}

func TestHandleErrorHTTPRequestError(t *testing.T) {
	requestError := http.NewHTTPRequestError(location+"TestHandleErrorHTTPRequestError", requestType, notification)
	result := model.HandleError(requestError)

	jsonRequestError, ok := result.(model.JSONRequestError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, requestError.Location, jsonRequestError.Location, test.EqualMessage)
	assert.Equal(t, requestError.RequestType, jsonRequestError.RequestType, test.EqualMessage)
	assert.Equal(t, requestError.Notification, jsonRequestError.Notification, test.EqualMessage)
}


