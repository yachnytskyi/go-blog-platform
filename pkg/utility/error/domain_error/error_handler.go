package domain_error

const (
	InternalErrorNotification       = "something went wrong, please repeat later"
	EntityNotFoundErrorNotification = "please repeat it later"
)

func ErrorHandler(err error) error {
	var errorMessage *ErrorMessage = new(ErrorMessage)

	switch err.(type) {
	case *ValidationError:
		return errorMessage
	case *EntityNotFoundError:
		errorMessage.Notification = EntityNotFoundErrorNotification
		return errorMessage
	default:
		errorMessage.Notification = InternalErrorNotification
		return errorMessage
	}
}

func ErrorsHandler(errors []error) []error {
	var errorMessage *ErrorMessage = new(ErrorMessage)

	for index, errorType := range errors {
		if _, ok := errorType.(*ValidationError); ok {
			continue

		} else if _, ok := errorType.(*EntityNotFoundError); ok {
			errorMessage.Notification = EntityNotFoundErrorNotification
			errors[index] = errorMessage
			return errors

		} else {
			errorMessage.Notification = InternalErrorNotification
			errors[index] = errorMessage
			return errors
		}
	}

	return errors
}
