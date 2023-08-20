package domain_error

const (
	InternalErrorNotification       = "something went wrong, please repeat later"
	EntityNotFoundErrorNotification = "please repeat it later"
)

func ErrorHandler(err error) error {
	switch errorType := err.(type) {
	case *ValidationError:
		return errorType
	case *ValidationErrors:
		return errorType
	case *EntityNotFoundError:
		var errorMessage *ErrorMessage = new(ErrorMessage)
		errorMessage.Notification = EntityNotFoundErrorNotification
		return errorMessage
	default:
		var errorMessage *ErrorMessage = new(ErrorMessage)
		errorMessage.Notification = InternalErrorNotification
		return errorMessage
	}
}
