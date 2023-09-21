package domain

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func IsUserOwner(currentUserID string, userID string) error {
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return domainError.NewErrorMessage("sorry, but you do not have permissions to do that")
	}
	return nil
}
