package domain

import (
	"fmt"

	"github.com/yachnytskyi/golang-mongo-grpc/config"
)

func HandleError(err error) error {
	fmt.Printf("Underlying Type: %T\n", err)
	fmt.Printf("Underlying Value: %v\n", err)
	switch errorType := err.(type) {
	case *ValidationError:
		return errorType
	case *ValidationErrors:
		return errorType
	case *EntityNotFoundError:
		return NewErrorMessage(config.EntityNotFoundErrorNotification)
	default:
		return NewErrorMessage(config.InternalErrorNotification)
	}
}
