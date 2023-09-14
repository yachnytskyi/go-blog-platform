package usecase

import (
	"fmt"
	"net"
	"strings"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
)

const (
	// Regex Patterns.
	emailRegex    = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\\\/=\\\\{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	usernameRegex = `^[a-zA-z0-9-_ \t]*$`
	passwordRegex = `^[a-zA-z0-9-_*,.]*$`
	minLength     = 4
	maxLength     = 40

	// Error Messages.
	stringAllowedLength       = "can be between %d and %d characters long"
	usernameAllowedCharacters = "sorry, only letters (a-z), numbers(0-9) and spaces are allowed"
	passwordAllowedCharacters = "sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed"
	emailAllowedCharacters    = "sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row"
	invalidEmailDomain        = "email domain does not exist"
	passwordsDoNotMatch       = "passwords do not match"

	// Field Names.
	nameField     = "name"
	EmailField    = "email"
	passwordField = "password"
	TypeRequired  = "required"
)

func validateUserCreate(userCreate userModel.UserCreate) (userModel.UserCreate, error) {
	validationErrors := &domainError.ValidationErrors{}
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	validateFieldError := validateField(userCreate.Name, nameField, TypeRequired, usernameRegex, usernameAllowedCharacters)
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
		return userModel.UserCreate{}, validationErrors
	}
	return userCreate, nil
}

func UserUpdateValidator(userUpdate userModel.UserUpdate) (userModel.UserUpdate, error) {
	validationErrors := &domainError.ValidationErrors{}
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)

	validateFieldError := validateField(userUpdate.Name, nameField, TypeRequired, usernameRegex, usernameAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return userModel.UserUpdate{}, validationErrors
	}
	return userUpdate, nil
}

func UserLoginValidator(userLogin userModel.UserLogin) (userModel.UserLogin, error) {
	validationErrors := &domainError.ValidationErrors{}
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)

	validateFieldError := validateEmail(userLogin.Email, EmailField, TypeRequired, emailRegex, emailAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	validateFieldError = validateField(userLogin.Password, passwordField, TypeRequired, usernameRegex, passwordAllowedCharacters)
	if validator.IsValueNotNil(validateFieldError) {
		validationErrors.ValidationErrors = append(validationErrors.ValidationErrors, validateFieldError)
	}
	if validator.IsSliceNotEmpty(validationErrors.ValidationErrors) {
		return userModel.UserLogin{}, validationErrors
	}
	return userLogin, nil
}

// func UserForgottenPasswordValidator(userForgottenPassword *userModel.UserForgottenPassword) error {
// 	var message string
// 	var err error

// 	if userForgottenPassword.Email == "" {
// 		message = "key: `UserForgottenPasswordView.Email` error: field validation for `email` failed, `email` cannot be empty "
// 		err = fmt.Errorf(message)
// 	}

// 	if validator.IsStringNotEmpty(userForgottenPassword.Email) {
// 		_, ok := mail.ParseAddress(userForgottenPassword.Email)
// 		if validator.IsValueNotNil(ok) {
// 			message = message + "key: `UserForgottenPasswordView.Email` error: field validation for `email` failed, invalid email address "
// 			err = fmt.Errorf(message)
// 		}
// 	}

// 	if validator.IsValueNotNil(err) {
// 		message = strings.TrimSpace(message)
// 		err = fmt.Errorf(message)

// 		return err
// 	}

// 	return nil
// }

// func UserResetPasswordValidator(userResetPassword *userModel.UserResetPassword) error {
// 	var message string
// 	var err error

// 	if userResetPassword.Password == "" {
// 		message = "key: `UserResetPasswordView.Password` error: field validation for `password` failed, `password` cannot be empty "
// 		err = fmt.Errorf(message)
// 	}

// 	if userResetPassword.PasswordConfirm == "" {
// 		message = "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, `password_confirm` cannot be empty "
// 		err = fmt.Errorf(message)
// 	}

// 	if userResetPassword.Password != "" && userResetPassword.PasswordConfirm != "" && userResetPassword.Password != userResetPassword.PasswordConfirm {
// 		message = message + "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match "
// 		err = fmt.Errorf(message)
// 	}

// 	if err != nil {
// 		message = strings.TrimSpace(message)
// 		err = fmt.Errorf(message)

// 		return err
// 	}

// 	return nil
// }

// func UserValidator(checkedStringKey string, checkedStringValue string, jsonStringKey string) string {
// 	var message string

// 	if checkedStringValue == "" {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" + " cannot be empty "
// 	}

// 	if len(checkedStringValue) > 100 {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" + " cannot be more than 100 characters long "
// 	}

// 	if !regexp.MustCompile(`^[a-zA-z0-9- \t]*$`).MatchString(checkedStringValue) {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
// 			" can use only letters, numbers, spaces, the hyphen or underscore character "
// 	}

// 	return message
// }

// func UserPasswordValidator(checkedStringKey string, checkedStringValue string, jsonStringKey string) string {
// 	var message string

// 	if checkedStringValue == "" {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
// 			" cannot be empty "
// 	}

// 	if len(checkedStringValue) > 100 {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
// 			" cannot be more than 100 characters long "
// 	}

// 	if !regexp.MustCompile(`^[a-zA-z0-9-]*$`).MatchString(checkedStringValue) {
// 		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
// 			" can use only letters, numbers, the hyphen or underscore character "
// 	}

// 	return message
// }

func IsEmailDomainValid(emailString string) bool {
	host := strings.Split(emailString, "@")[1]
	_, err := net.LookupMX(host)
	if validator.IsErrorNotNil(err) {
		return false
	}
	return true
}

func validateField(field, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(field, minLength, maxLength)) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(stringAllowedLength, minLength, maxLength))
	} else if domainUtility.CheckSpecialCharactersString(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	}
	return domainError.ValidationError{}
}
func validateEmail(field, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(field, minLength, maxLength)) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(stringAllowedLength, minLength, maxLength))
	} else if domainUtility.CheckSpecialCharactersString(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	} else if validator.IsBooleanNotTrue(IsEmailDomainValid(field)) {
		return domainError.NewValidationError(fieldName, fieldType, invalidEmailDomain)
	}
	return domainError.ValidationError{}
}
func validatePassword(password, passwordConfirm, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(password, minLength, maxLength)) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(stringAllowedLength, minLength, maxLength))
	} else if domainUtility.CheckSpecialCharactersString(password, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	} else if validator.AreStringsNotEqual(password, passwordConfirm) {
		return domainError.NewValidationError(fieldName, fieldType, passwordsDoNotMatch)
	}
	return domainError.ValidationError{}
}
