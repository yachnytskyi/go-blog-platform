package usecase

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
)

const (
	emailRegex                = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\\\/=\\\\{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	usernameRegex             = `^[a-zA-z0-9-_ \t]*$`
	passwordRegex             = `^[a-zA-z0-9-_*,.]*$`
	minLength                 = 4
	maxLength                 = 40
	usernameAllowedCharacters = "sorry, only letters (a-z), numbers(0-9) and spaces are allowed"
	passwordAllowedCharacters = "sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed"
	invalidEmail              = "sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row"
	invalidEmailDomain        = "email domain does not exist"
	invalidPassword           = "passwords do not match"
	nameField                 = "name"
	emailField                = "email"
	passwordField             = "password"
	typeRequired              = "required"
	typeOptional              = "optional"
)

func UserCreateValidator(userCreate *userModel.UserCreate) error {
	stringAllowedLength := "can be between " + strconv.Itoa(minLength) + " and " + strconv.Itoa(maxLength)
	userCreateValidationErrors := &domainError.ValidationErrors{}
	domainUtility.SanitizeString(&userCreate.Name)
	domainUtility.SanitizeString(&userCreate.Email)
	domainUtility.SanitizeString(&userCreate.Password)
	domainUtility.SanitizeString(&userCreate.PasswordConfirm)
	domainUtility.StringToLower(&userCreate.Email)

	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(userCreate.Name, minLength, maxLength)) {
		userCreateValidationError := domainError.NewValidationError(nameField, typeRequired, nameField+" "+stringAllowedLength)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	} else if domainUtility.CheckSpecialCharactersString(userCreate.Name, usernameRegex) {
		userCreateValidationError := domainError.NewValidationError(nameField, typeRequired, usernameAllowedCharacters)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	}
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(userCreate.Email, minLength, maxLength)) {
		userCreateValidationError := domainError.NewValidationError(emailField, typeRequired, emailField+" "+stringAllowedLength)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	} else if domainUtility.CheckSpecialCharactersString(userCreate.Email, emailRegex) {
		userCreateValidationError := domainError.NewValidationError(emailField, typeRequired, invalidEmail)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	} else if validator.IsBooleanNotTrue(IsEmailDomainValid(userCreate.Email)) {
		userCreateValidationError := domainError.NewValidationError(emailField, typeRequired, invalidEmailDomain)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	}
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(userCreate.Password, minLength, maxLength)) {
		userCreateValidationError := domainError.NewValidationError(passwordField, typeRequired, passwordField+" "+stringAllowedLength)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	} else if validator.IsBooleanNotTrue(validator.CheckMatchStrings(userCreate.Password, userCreate.PasswordConfirm)) {
		userCreateValidationError := domainError.NewValidationError(passwordField, typeRequired, invalidPassword)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	} else if domainUtility.CheckSpecialCharactersString(userCreate.Password, passwordRegex) {
		userCreateValidationError := domainError.NewValidationError(passwordField, typeRequired, passwordAllowedCharacters)
		userCreateValidationErrors.ValidationErrors = append(userCreateValidationErrors.ValidationErrors, userCreateValidationError)
	}
	if validator.IsSliceNotEmpty(userCreateValidationErrors.ValidationErrors) {
		return userCreateValidationErrors
	}
	return nil
}

func UserUpdateValidator(userUpdate *userModel.UserUpdate) error {
	userUpdateValidationErrors := &domainError.ValidationErrors{}
	stringAllowedLength := "can be between " + strconv.Itoa(minLength) + " and " + strconv.Itoa(maxLength)
	domainUtility.SanitizeString(&userUpdate.Name)
	if validator.IsBooleanNotTrue(domainUtility.CheckCorrectLengthString(userUpdate.Name, minLength, maxLength)) {
		userUpdateValidationError := domainError.ValidationError{
			Field:        "name",
			FieldType:    "required",
			Notification: stringAllowedLength,
		}
		userUpdateValidationErrors.ValidationErrors = append(userUpdateValidationErrors.ValidationErrors, userUpdateValidationError)

	} else if domainUtility.CheckSpecialCharactersString(userUpdate.Name, usernameRegex) {
		userUpdateValidationError := domainError.ValidationError{
			Field:        "name",
			FieldType:    "required",
			Notification: usernameAllowedCharacters,
		}
		userUpdateValidationErrors.ValidationErrors = append(userUpdateValidationErrors.ValidationErrors, userUpdateValidationError)
	}
	if validator.IsSliceNotEmpty(userUpdateValidationErrors.ValidationErrors) {
		return userUpdateValidationErrors
	}
	return nil
}

func UserLoginValidator(userLogin *userModel.UserLogin) error {
	var message string
	var err error

	if userLogin.Email == "" {
		message = "key: `UserLogInView.Email` error: field validation for `email` failed, `email` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userLogin.Password == "" {
		message = message + "key: `UserLogIn.Password` error: field validation for `password` failed, `password` cannot be empty "
		err = fmt.Errorf(message)
	}

	if validator.IsValueNotNil(err) {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)
		return err
	}
	return nil
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
