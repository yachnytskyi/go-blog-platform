package model

import (
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HTTPAuthorizationErrorToJSONAuthorizationErrorMapper(httpAuthorizationError httpError.HTTPAuthorizationError) JSONAuthorizationError {
	return NewJSONAuthorizationError(
		httpAuthorizationError.Location,
		httpAuthorizationError.Notification,
	)
}

func HTTPRequestErrorToJSONRequestErrorMapper(httpRequestError httpError.HTTPRequestError) JSONRequestError {
	return NewJSONRequestError(
		httpRequestError.Location,
		httpRequestError.RequestType,
		httpRequestError.Notification,
	)
}

func HTTPInternalErrorToJSONInternalErrorMapper(httpInternalError httpError.HTTPInternalError) JSONInternalError {
	return NewJSONInternalError(
		httpInternalError.Location,
		httpInternalError.Notification,
	)
}

func InfoMessageToJSONInfoMessageMapper(infoMessage domainError.InfoMessage) JSONInfoMessage {
	return NewJSONInfoMessage(
		infoMessage.Location,
		infoMessage.Notification,
	)
}

func ValidationErrorToJSONValidationErrorMapper(validationError domainError.ValidationError) JSONValidationError {
	return NewJSONValidationError(
		validationError.Location,
		validationError.Field,
		validationError.FieldType,
		validationError.Notification,
	)
}

func ValidationErrorsToJSONValidationErrorsMapper(validationErrors domainError.ValidationErrors) JSONValidationErrors {
	JSONValidationErrors := make([]error, 0, validationErrors.Len())
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domainError.ValidationError)
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

func AuthorizationErrorToJSONAuthorizationErrorMapper(authorizationError domainError.AuthorizationError) JSONAuthorizationError {
	return NewJSONAuthorizationError(
		authorizationError.Location,
		authorizationError.Notification,
	)
}

func ItemNotFoundErrorToJSONItemNotFoundErrorMapper(itemNotFoundError domainError.ItemNotFoundError) JSONItemNotFoundError {
	return NewJSONItemNotFoundError(
		itemNotFoundError.Location,
		itemNotFoundError.Notification,
		itemNotFoundError.Query,
	)
}

func InvalidTokenErrorToJSONIvalidTokenErrorMapper(invalidTokenError domainError.InvalidTokenError) JSONInvalidTokenError {
	return NewJSONInvalidTokenError(
		invalidTokenError.Location,
		invalidTokenError.Notification,
	)
}

func TimeExpiredErrorToJSONTimeExpiredErrorMapper(timeExpiredError domainError.TimeExpiredError) JSONTimeExpiredError {
	return NewJSONTimeExpiredError(
		timeExpiredError.Location,
		timeExpiredError.Notification,
	)
}

func PaginationErrorToJSONPaginationErrorMapper(paginationError domainError.PaginationError) JSONPaginationError {
	return NewJSONPaginationError(
		paginationError.Location,
		paginationError.CurrentPage,
		paginationError.TotalPages,
		paginationError.Notification,
	)
}

func InternalErrorToJSONInternalErrorMapper(internalError domainError.InternalError) JSONInternalError {
	return NewJSONInternalError(
		internalError.Location,
		internalError.Notification,
	)
}
