package domain

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func IsUserOwner(currentUserID string, userID string) error {
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return domainError.NewAuthorizationError(location+"IsUserOwner.AreStringsNotEqual", constants.AuthorizationErrorNotification)
	}

	return nil
}
