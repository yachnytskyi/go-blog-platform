package utility

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(location, password string) (string, error) {
	hashedPassword, generateFromPasswordError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if validator.IsError(generateFromPasswordError) {
		internalError := domainError.NewInternalError(location+".HashPassword.bcrypt.GenerateFromPassword", generateFromPasswordError.Error())
		logging.Logger(internalError)
		return constants.EmptyString, internalError
	}
	return string(hashedPassword), nil
}
