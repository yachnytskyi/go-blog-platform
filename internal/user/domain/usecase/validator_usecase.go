package usecase

import (
	"fmt"
	"net"
	"strings"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
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
	usernameAllowedCharacters = "sorry, only letters (a-z), numbers(0-9) and spaces are allowed"
	passwordAllowedCharacters = "sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed"
	emailAllowedCharacters    = "sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row"
	invalidEmailDomain        = "email domain does not exist"
	passwordsDoNotMatch       = "passwords do not match"
	invalidEmailOrPassword    = "invalid email or password"

	// Field Names.
	nameField     = "name"
	EmailField    = "email"
	passwordField = "password"
	TypeRequired  = "required"
)

func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	validationErrors := domainError.ValidationErrors{}
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validateFieldError := domainUtility.ValidateField(userCreate.Name, nameField, TypeRequired, usernameRegex, usernameAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = validateEmail(userCreate.Email, EmailField, TypeRequired, emailRegex, emailAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = validatePassword(userCreate.Password, userCreate.PasswordConfirm, passwordField, TypeRequired, usernameRegex, passwordAllowedCharacters)
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

	validateFieldError := domainUtility.ValidateField(userUpdate.Name, nameField, TypeRequired, usernameRegex, usernameAllowedCharacters)
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

	validateFieldError := validateEmail(userLogin.Email, EmailField, TypeRequired, emailRegex, emailAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = domainUtility.ValidateField(userLogin.Password, passwordField, TypeRequired, usernameRegex, passwordAllowedCharacters)
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

	validateFieldError := validateEmail(userForgottenPassword.Email, EmailField, TypeRequired, emailRegex, emailAllowedCharacters)
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

	validateFieldError := validatePassword(userResetPassword.Password, userResetPassword.PasswordConfirm, passwordField, TypeRequired, usernameRegex, passwordAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](validationErrors)
	}
	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

func validateEmail(field, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if domainUtility.IsStringLengthNotCorrect(field, config.MinLength, config.MaxLength) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(config.StringAllowedLength, config.MinLength, config.MaxLength))
	} else if domainUtility.CheckSpecialCharactersString(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	} else if isEmailDomainNotValid(field) {
		return domainError.NewValidationError(fieldName, fieldType, invalidEmailDomain)
	}
	return domainError.ValidationError{}
}

func validatePassword(password, passwordConfirm, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if domainUtility.IsStringLengthNotCorrect(password, config.MinLength, config.MaxLength) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(config.StringAllowedLength, config.MinLength, config.MaxLength))
	} else if domainUtility.CheckSpecialCharactersString(password, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	} else if validator.AreStringsNotEqual(password, passwordConfirm) {
		return domainError.NewValidationError(fieldName, fieldType, passwordsDoNotMatch)
	}
	return domainError.ValidationError{}
}

func isEmailDomainNotValid(emailString string) bool {
	host := strings.Split(emailString, "@")[1]
	_, err := net.LookupMX(host)
	if validator.IsErrorNotNil(err) {
		return true
	}
	return false
}

// Compare the encrypted and the user provided passwords.
func arePasswordsEqual(hashedPassword string, checkedPassword string) bool {
	if validator.IsErrorNotNil(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		return false
	}
	return true
}

// Compare the encrypted and the user provided passwords.
func arePasswordsNotEqual(hashedPassword string, checkedPassword string) bool {
	if validator.IsErrorNotNil(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		return true
	}
	return false
}
