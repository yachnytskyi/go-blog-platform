package usecase

import (
	"fmt"
	"net"
	"strings"

	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"

	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
	bcrypt "golang.org/x/crypto/bcrypt"
)

const (
	// Regex Patterns.
	emailRegex    = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\\\/=\\\\{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	usernameRegex = `^[a-zA-z0-9-_ \t]*$`
	passwordRegex = `^[a-zA-z0-9-_*,.]*$`

	// Error Messages.
	usernameAllowedCharacters = "sorry, only letters (a-z), numbers(0-9) and spaces are allowed."
	passwordAllowedCharacters = "sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed."
	emailAllowedCharacters    = "sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row."
	invalidEmailDomain        = "email domain does not exist."
	passwordsDoNotMatch       = "passwords do not match."
	invalidEmailOrPassword    = "invalid email or password."

	// Field Names.
	nameField     = "name"
	EmailField    = "email"
	passwordField = "password"

	// Location.
	location = "internal.user.domain.usecase."
)

func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	validationErrors := domainError.ValidationErrors{}
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validateFieldError := validateEmail(userCreate.Email, emailRegex)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = domainUtility.ValidateField(userCreate.Name, nameField, usernameRegex, usernameAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = validatePassword(userCreate.Password, userCreate.PasswordConfirm, passwordField, usernameRegex)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserCreate](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserCreate](userCreate)
}

func validateUserUpdate(userUpdate userModel.UserUpdate) common.Result[userModel.UserUpdate] {
	validationErrors := domainError.ValidationErrors{}
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)

	validateFieldError := domainUtility.ValidateField(userUpdate.Name, nameField, usernameRegex, usernameAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserUpdate](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserUpdate](userUpdate)
}

func validateUserLogin(userLogin userModel.UserLogin) common.Result[userModel.UserLogin] {
	validationErrors := domainError.ValidationErrors{}
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)

	validateFieldError := validateEmail(userLogin.Email, emailRegex)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = domainUtility.ValidateField(userLogin.Password, passwordField, passwordRegex, passwordAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserLogin](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserLogin](userLogin)
}

func validateUserForgottenPassword(userForgottenPassword userModel.UserForgottenPassword) common.Result[userModel.UserForgottenPassword] {
	validationErrors := domainError.ValidationErrors{}
	userForgottenPassword.Email = domainUtility.SanitizeAndToLowerString(userForgottenPassword.Email)

	validateFieldError := validateEmail(userForgottenPassword.Email, emailRegex)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserForgottenPassword](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserForgottenPassword](userForgottenPassword)
}

func validateResetPassword(userResetPassword userModel.UserResetPassword) common.Result[userModel.UserResetPassword] {
	validationErrors := domainError.ValidationErrors{}
	userResetPassword.Password = domainUtility.SanitizeString(userResetPassword.Password)
	userResetPassword.PasswordConfirm = domainUtility.SanitizeString(userResetPassword.PasswordConfirm)

	validateFieldError := validatePassword(userResetPassword.Password, userResetPassword.PasswordConfirm, passwordField, passwordRegex)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

func validateEmail(email, fieldRegex string) domainError.ValidationError {
	if domainUtility.IsStringLengthNotValid(email, constant.MinStringLength, constant.MaxStringLength) {
		notification := fmt.Sprintf(constant.StringAllowedLength, constant.MinStringLength, constant.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthNotValid", EmailField, constant.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if domainUtility.IsStringCharactersNotValid(email, fieldRegex) {
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringCharactersNotValid", EmailField, constant.FieldRequired, emailAllowedCharacters)
		logging.Logger(validationError)
		return validationError
	}
	if isEmailDomainNotValid(email) {
		validationError := domainError.NewValidationError(location+"validateEmail.IsEmailDomainNotValid", EmailField, constant.FieldRequired, invalidEmailDomain)
		logging.Logger(validationError)
		return validationError
	}
	return domainError.ValidationError{}
}

func validatePassword(password, passwordConfirm, fieldName, fieldRegex string) domainError.ValidationError {
	if domainUtility.IsStringLengthNotValid(password, constant.MinStringLength, constant.MaxStringLength) {
		notification := fmt.Sprintf(constant.StringAllowedLength, constant.MinStringLength, constant.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validatePassword.IsStringLengthNotValid", fieldName, constant.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if domainUtility.IsStringCharactersNotValid(password, fieldRegex) {
		validationError := domainError.NewValidationError(location+"validatePassword.IsStringCharactersNotValid", fieldName, constant.FieldRequired, passwordAllowedCharacters)
		logging.Logger(validationError)
		return validationError
	}
	if validator.AreStringsNotEqual(password, passwordConfirm) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringsNotEqual", fieldName, constant.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		return validationError
	}
	return domainError.ValidationError{}
}

func isEmailDomainNotValid(emailString string) bool {
	host := strings.Split(emailString, "@")[1]
	_, lookupMXError := net.LookupMX(host)
	return validator.IsErrorNotNil(lookupMXError)
}

// Compare the encrypted and the user provided passwords.
func arePasswordsNotEqual(hashedPassword string, checkedPassword string) domainError.ValidationError {
	if validator.IsErrorNotNil(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		validationError := domainError.NewValidationError(location+"arePasswordsNotEqual.CompareHashAndPassword", checkedPassword, constant.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		return validationError
	}
	return domainError.ValidationError{}
}
