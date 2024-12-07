package model

import (
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HTTPAuthorizationErrorToJSONAuthorizationErrorMapper(httpAuthorizationError http.HTTPAuthorizationError) JSONAuthorizationError {
	return NewJSONAuthorizationError(
		httpAuthorizationError.Location,
		httpAuthorizationError.Notification,
	)
}

func HTTPRequestErrorToJSONRequestErrorMapper(httpRequestError http.HTTPRequestError) JSONRequestError {
	return NewJSONRequestError(
		httpRequestError.Location,
		httpRequestError.RequestType,
		httpRequestError.Notification,
	)
}

func HTTPInternalErrorToJSONInternalErrorMapper(httpInternalError http.HTTPInternalError) JSONInternalError {
	return NewJSONInternalError(
		httpInternalError.Location,
		httpInternalError.Notification,
	)
}

func InfoMessageToJSONInfoMessageMapper(infoMessage domain.InfoMessage) JSONInfoMessage {
	return NewJSONInfoMessage(
		infoMessage.Location,
		infoMessage.Notification,
	)
}

func ValidationErrorToJSONValidationErrorMapper(validationError domain.ValidationError) JSONValidationError {
	return NewJSONValidationError(
		validationError.Location,
		validationError.Field,
		validationError.FieldType,
		validationError.Notification,
	)
}

func ValidationErrorsToJSONValidationErrorsMapper(validationErrors domain.ValidationErrors) JSONValidationErrors {
	JSONValidationErrors := make([]error, 0, len(validationErrors.Errors))
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domain.ValidationError)
		if ok {
			JSONValidationError := NewJSONValidationError(
				validationError.Location,
				validationError.Field,
				validationError.FieldType,
				validationError.Notification,
			)
			JSONValidationErrors = append(JSONValidationErrors, JSONValidationError)
		}
	}

	return NewJSONValidationErrors(JSONValidationErrors)
}

func AuthorizationErrorToJSONAuthorizationErrorMapper(authorizationError domain.AuthorizationError) JSONAuthorizationError {
	return NewJSONAuthorizationError(
		authorizationError.Location,
		authorizationError.Notification,
	)
}

func ItemNotFoundErrorToJSONItemNotFoundErrorMapper(itemNotFoundError domain.ItemNotFoundError) JSONItemNotFoundError {
	return NewJSONItemNotFoundError(
		itemNotFoundError.Location,
		itemNotFoundError.Notification,
		itemNotFoundError.Query,
	)
}

func InvalidTokenErrorToJSONIvalidTokenErrorMapper(invalidTokenError domain.InvalidTokenError) JSONInvalidTokenError {
	return NewJSONInvalidTokenError(
		invalidTokenError.Location,
		invalidTokenError.Notification,
	)
}

func TimeExpiredErrorToJSONTimeExpiredErrorMapper(timeExpiredError domain.TimeExpiredError) JSONTimeExpiredError {
	return NewJSONTimeExpiredError(
		timeExpiredError.Location,
		timeExpiredError.Notification,
	)
}

func PaginationErrorToJSONPaginationErrorMapper(paginationError domain.PaginationError) JSONPaginationError {
	return NewJSONPaginationError(
		paginationError.Location,
		paginationError.CurrentPage,
		paginationError.TotalPages,
		paginationError.Notification,
	)
}

func InternalErrorToJSONInternalErrorMapper(internalError domain.InternalError) JSONInternalError {
	return NewJSONInternalError(
		internalError.Location,
		internalError.Notification,
	)
}
