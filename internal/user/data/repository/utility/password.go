package utility

import (
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hashed password from a plain text password.
func HashPassword(location, password string) commonModel.Result[string] {
	hashedPassword, generateFromPasswordError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if validator.IsError(generateFromPasswordError) {
		internalError := domainError.NewInternalError(location+".HashPassword.GenerateFromPassword", generateFromPasswordError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[string](internalError)
	}

	return commonModel.NewResultOnSuccess(string(hashedPassword))
}
