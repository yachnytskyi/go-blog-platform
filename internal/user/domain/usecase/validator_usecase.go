package usecase

import (
	"fmt"
	"net"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
	bcrypt "golang.org/x/crypto/bcrypt"
)

const (
	// Location.
	location = "internal.user.domain.usecase."

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
	nameField             = "name"
	EmailField            = "email"
	passwordField         = "password"
	emailOrPasswordFields = "email or password"

	// Amount of expected possinle errors
	expectedErrors = 4
)

func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	validationErrors := make([]error, 0, expectedErrors)
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validateEmailError := validateEmail(userCreate.Email, emailRegex)
	if validator.IsError(validateEmailError) {
		validationErrors = append(validationErrors, validateEmailError)
	}
	validateNameError := domainUtility.ValidateField(userCreate.Name, nameField, usernameRegex, usernameAllowedCharacters)
	if validator.IsError(validateNameError) {
		validationErrors = append(validationErrors, validateNameError)
	}
	validatePasswordError := validatePassword(userCreate.Password, userCreate.PasswordConfirm, passwordField, usernameRegex)
	if validator.IsError(validatePasswordError) {
		validationErrors = append(validationErrors, validatePasswordError)
	}
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserCreate](domainError.ValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserCreate](userCreate)
}

func validateUserUpdate(userUpdate userModel.UserUpdate) common.Result[userModel.UserUpdate] {
	validationErrors := make([]error, 0, expectedErrors)
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)
	validateFieldError := domainUtility.ValidateField(userUpdate.Name, nameField, usernameRegex, usernameAllowedCharacters)
	if validator.IsError(validateFieldError) {
		validationErrors = append(validationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserUpdate](domainError.ValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserUpdate](userUpdate)
}

func validateUserLogin(userLogin userModel.UserLogin) common.Result[userModel.UserLogin] {
	validationErrors := make([]error, 0, expectedErrors)
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)
	validateFieldError := validateEmail(userLogin.Email, emailRegex)
	if validator.IsError(validateFieldError) {
		validationErrors = append(validationErrors, validateFieldError)
	}
	validateFieldError = domainUtility.ValidateField(userLogin.Password, passwordField, passwordRegex, passwordAllowedCharacters)
	if validator.IsError(validateFieldError) {
		validationErrors = append(validationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserLogin](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserLogin](userLogin)
}

func validateUserForgottenPassword(userForgottenPassword userModel.UserForgottenPassword) common.Result[userModel.UserForgottenPassword] {
	var validationErrors domainError.ValidationErrors
	userForgottenPassword.Email = domainUtility.SanitizeAndToLowerString(userForgottenPassword.Email)
	validateFieldError := validateEmail(userForgottenPassword.Email, emailRegex)
	if validator.IsError(validateFieldError) {
		validationErrors = append(validationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserForgottenPassword](validationErrors)
	}

	return common.NewResultOnSuccess[userModel.UserForgottenPassword](userForgottenPassword)
}

func validateResetPassword(userResetPassword userModel.UserResetPassword) common.Result[userModel.UserResetPassword] {
	var validationErrors domainError.ValidationErrors
	userResetPassword.Password = domainUtility.SanitizeString(userResetPassword.Password)
	userResetPassword.PasswordConfirm = domainUtility.SanitizeString(userResetPassword.PasswordConfirm)

	validateFieldError := validatePassword(userResetPassword.Password, userResetPassword.PasswordConfirm, passwordField, passwordRegex)
	if validator.IsError(validateFieldError) {
		validationErrors = append(validationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](validationErrors)
	}

	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

func validateEmail(email, fieldRegex string) error {
	if domainUtility.IsStringLengthInvalid(email, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthInvalid", EmailField, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if domainUtility.AreStringCharactersInvalid(email, fieldRegex) {
		validationError := domainError.NewValidationError(location+"validateEmail.AreStringCharactersInvalid", EmailField, constants.FieldRequired, emailAllowedCharacters)
		logging.Logger(validationError)
		return validationError
	}
	if isEmailDomainNotValid(email) {
		validationError := domainError.NewValidationError(location+"validateEmail.IsEmailDomainNotValid", EmailField, constants.FieldRequired, invalidEmailDomain)
		logging.Logger(validationError)
		return validationError
	}

	return nil
}

func validatePassword(password, passwordConfirm, fieldName, fieldRegex string) error {
	if domainUtility.IsStringLengthInvalid(password, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validatePassword.IsStringLengthInvalid", fieldName, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if domainUtility.AreStringCharactersInvalid(password, fieldRegex) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringCharactersInvalid", fieldName, constants.FieldRequired, passwordAllowedCharacters)
		logging.Logger(validationError)
		return validationError
	}
	if validator.AreStringsNotEqual(password, passwordConfirm) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringsNotEqual", fieldName, constants.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		return validationError
	}

	return nil
}

func isEmailDomainNotValid(emailString string) bool {
	host := strings.Split(emailString, "@")[1]
	_, lookupMXError := net.LookupMX(host)
	return validator.IsError(lookupMXError)
}

// Compare the encrypted and the user provided passwords.
func arePasswordsNotEqual(hashedPassword string, checkedPassword string) error {
	if validator.IsError(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		validationError := domainError.NewValidationError(location+"arePasswordsNotEqual.CompareHashAndPassword", emailOrPasswordFields, constants.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		validationError.Notification = invalidEmailOrPassword
		return validationError
	}

	return nil
}
