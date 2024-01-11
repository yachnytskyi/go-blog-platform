package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func ValidationErrorToHttpValidationErrorViewMapper(validationError domainError.ValidationError) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:        validationError.Field,
		FieldType:    validationError.FieldType,
		Notification: validationError.Notification,
	}
}

func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors domainError.ValidationErrors) HttpValidationErrorsView {
	httpValidationErrorsView := make([]HttpValidationErrorView, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		validationError, ok := validationError.(domainError.ValidationError)
		if ok {
			// Map the specific validation error to the HTTP view.
			httpValidationErrorView := NewHttpValidationError(validationError.Field, validationError.FieldType, validationError.Notification)
			httpValidationErrorsView = append(httpValidationErrorsView, httpValidationErrorView)
		}
	}
	return HttpValidationErrorsView(httpValidationErrorsView)
}

func AuthorizationErrorToHttpAuthorizationErrorViewMapper(authorizationError domainError.AuthorizationError) HttpAuthorizationErrorView {
	return HttpAuthorizationErrorView{
		Notification: authorizationError.Notification,
	}
}

func EntityNotFoundErrorToHttpEntityNotFoundErrorViewMapper(entityNotFoundError domainError.EntityNotFoundError) HttpEntityNotFoundErrorView {
	return HttpEntityNotFoundErrorView{
		Notification: entityNotFoundError.Notification,
	}
}

func PaginationErrorToHttpPaginationErrorViewMapper(errorMessage domainError.PaginationError) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		CurrentPage:  errorMessage.CurrentPage,
		TotalPages:   errorMessage.TotalPages,
		Notification: errorMessage.Notification,
	}
}

func InternalErrorToHttpInternalErrorViewMapper(internalError domainError.InternalError) HttpInternalErrorView {
	return HttpInternalErrorView{
		Notification: internalError.Notification,
	}
}
