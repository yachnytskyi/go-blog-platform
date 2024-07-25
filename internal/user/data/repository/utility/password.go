package utility

import (
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hashed password from a plain text password.
func HashPassword(logger model.Logger, location, password string) common.Result[string] {
	hashedPassword, generateFromPasswordError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if validator.IsError(generateFromPasswordError) {
		internalError := domainError.NewInternalError(location+".HashPassword.GenerateFromPassword", generateFromPasswordError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[string](internalError)
	}

	return common.NewResultOnSuccess(string(hashedPassword))
}
