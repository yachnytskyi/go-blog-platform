package utility

import (
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hashed password from a plain text password.
func HashPassword(logger applicationModel.Logger, location, password string) commonModel.Result[string] {
	hashedPassword, generateFromPasswordError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if validator.IsError(generateFromPasswordError) {
		internalError := domainError.NewInternalError(location+".HashPassword.GenerateFromPassword", generateFromPasswordError.Error())
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[string](internalError)
	}

	return commonModel.NewResultOnSuccess(string(hashedPassword))
}
