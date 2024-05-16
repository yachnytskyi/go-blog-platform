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
	usernameField         = "name"
	EmailField            = "email"
	passwordField         = "password"
	emailOrPasswordFields = "email or password"

	// Amount of expected possinle errors
	expectedErrors = 4
)

var (
	emailValidator = domainUtility.CommonValidator{
		FieldName:  EmailField,
		FieldRegex: emailRegex,
		MinLength:  constants.MinStringLength,
		MaxLength:  constants.MaxStringLength,
	}
	usernameValidator = domainUtility.CommonValidator{
		FieldName:    usernameField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: usernameAllowedCharacters,
	}
	passwordValidator = domainUtility.CommonValidator{
		FieldName:  passwordField,
		FieldRegex: usernameRegex,
		MinLength:  constants.MinStringLength,
		MaxLength:  constants.MaxStringLength,
	}
	// Add more validators for other fields as needed
)

func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	validationErrors := make([]error, 0, expectedErrors)
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validationErrors = validateEmail(userCreate.Email, validationErrors)
	validationErrors = domainUtility.ValidateField(userCreate.Name, usernameValidator, validationErrors)
	validationErrors = validatePassword(userCreate.Password, userCreate.PasswordConfirm, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserCreate](domainError.ValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserCreate](userCreate)
}

func validateUserUpdate(userUpdate userModel.UserUpdate) common.Result[userModel.UserUpdate] {
	validationErrors := make([]error, expectedErrors)
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)

	validationErrors = domainUtility.ValidateField(userUpdate.Name, usernameValidator, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserUpdate](domainError.ValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserUpdate](userUpdate)
}

func validateUserLogin(userLogin userModel.UserLogin) common.Result[userModel.UserLogin] {
	validationErrors := make([]error, expectedErrors)
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)

	validationErrors = validateEmail(userLogin.Email, validationErrors)
	validationErrors = domainUtility.ValidateField(userLogin.Password, usernameValidator, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserLogin](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserLogin](userLogin)
}

func validateUserForgottenPassword(userForgottenPassword userModel.UserForgottenPassword) common.Result[userModel.UserForgottenPassword] {
	validationErrors := make([]error, expectedErrors)
	userForgottenPassword.Email = domainUtility.SanitizeAndToLowerString(userForgottenPassword.Email)
	validationErrors = validateEmail(userForgottenPassword.Email, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserForgottenPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserForgottenPassword](userForgottenPassword)
}

func validateResetPassword(userResetPassword userModel.UserResetPassword) common.Result[userModel.UserResetPassword] {
	validationErrors := make([]error, expectedErrors)
	userResetPassword.Password = domainUtility.SanitizeString(userResetPassword.Password)
	userResetPassword.PasswordConfirm = domainUtility.SanitizeString(userResetPassword.PasswordConfirm)

	validationErrors = validatePassword(userResetPassword.Password, userResetPassword.PasswordConfirm, validationErrors)
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

func validateEmail(email string, validationErrors []error) []error {
	// Preallocate enough capacity for all elements but set length to zero.
	// Append initial elements.
	errors := make([]error, len(validationErrors))
	errors = append(errors, validationErrors...)
	if domainUtility.IsStringLengthInvalid(email, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthInvalid", EmailField, constants.FieldRequired, notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}
	if domainUtility.AreStringCharactersInvalid(email, emailValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"validateEmail.AreStringCharactersInvalid", EmailField, constants.FieldRequired, emailAllowedCharacters)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}
	if isEmailDomainNotValid(email) {
		validationError := domainError.NewValidationError(location+"validateEmail.IsEmailDomainNotValid", EmailField, constants.FieldRequired, invalidEmailDomain)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

func validatePassword(password, passwordConfirm string, validationErrors []error) []error {
	// Preallocate enough capacity for all elements but set length to zero.
	// Append initial elements.
	errors := make([]error, len(validationErrors))
	errors = append(errors, validationErrors...)
	if domainUtility.IsStringLengthInvalid(password, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, passwordValidator.MinLength, passwordValidator.MaxLength)
		validationError := domainError.NewValidationError(location+"validatePassword.IsStringLengthInvalid", passwordValidator.FieldName, constants.FieldRequired, notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}
	if domainUtility.AreStringCharactersInvalid(password, passwordValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringCharactersInvalid", passwordValidator.FieldName, constants.FieldRequired, passwordAllowedCharacters)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}
	if validator.AreStringsNotEqual(password, passwordConfirm) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringsNotEqual", passwordValidator.FieldName, constants.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
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

func isEmailValid(email string) error {
	if domainUtility.IsStringLengthInvalid(email, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthInvalid", EmailField, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if domainUtility.AreStringCharactersInvalid(email, emailValidator.FieldRegex) {
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
