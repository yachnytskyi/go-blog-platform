package model

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	domainValidatorUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain_validator"
)

const (
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

func (userCreate *UserCreate) UserCreateValidator() []*domainError.DomainValidationError {
	var userErrors []*domainError.DomainValidationError

	emptyUsername := false
	emptyEmail := false
	emptyPassword := false
	validEmailAddress := true
	passwordsMatch := true

	domainValidatorUtility.SanitizeString(&userCreate.Name)
	domainValidatorUtility.SanitizeString(&userCreate.Email)
	domainValidatorUtility.SanitizeString(&userCreate.Password)
	domainValidatorUtility.SanitizeString(&userCreate.PasswordConfirm)

	if domainValidatorUtility.IsStringNull(userCreate.Name) {
		userError := &domainError.DomainValidationError{
			Field:        "name",
			FieldType:    "required",
			Notification: "cannot be empty",
		}

		emptyUsername = true
		userErrors = append(userErrors, userError)
	}

	if !emptyUsername {
		if domainValidatorUtility.IsStringLengthExceeded(userCreate.Name) {
			userError := &domainError.DomainValidationError{
				Field:        "name",
				FieldType:    "required",
				Notification: "cannot be more than 100 characters long",
			}

			userErrors = append(userErrors, userError)
		}

		if IsUserNameContainsSpecialCharacters(userCreate.Name) {
			userError := &domainError.DomainValidationError{
				Field:        "name",
				FieldType:    "required",
				Notification: "sorry, only letters (a-z), numbers(0-9) and spaces are allowed",
			}
			userErrors = append(userErrors, userError)
		}
	}

	if domainValidatorUtility.IsStringNull(userCreate.Email) {
		userError := &domainError.DomainValidationError{
			Field:        "email",
			FieldType:    "required",
			Notification: "cannot be empty",
		}

		emptyEmail = true
		userErrors = append(userErrors, userError)
	}

	if !emptyEmail {
		if domainValidatorUtility.IsStringLengthExceeded(userCreate.Email) {
			userError := &domainError.DomainValidationError{
				Field:        "email",
				FieldType:    "required",
				Notification: "cannot be more than 100 characters long",
			}

			userErrors = append(userErrors, userError)
		}

		if !IsEmailValid(userCreate.Email) {
			userError := &domainError.DomainValidationError{
				Field:        "email",
				FieldType:    "required",
				Notification: "invalid email address",
			}

			validEmailAddress = false
			userErrors = append(userErrors, userError)
		}
	}

	if validEmailAddress {
		if !IsEmailDomainValid(userCreate.Email) {
			userError := &domainError.DomainValidationError{
				Field:        "email",
				FieldType:    "required",
				Notification: "email domain does not exist",
			}

			userErrors = append(userErrors, userError)
		}
	}

	if domainValidatorUtility.IsStringNull(userCreate.Password) {
		userError := &domainError.DomainValidationError{
			Field:        "password",
			FieldType:    "required",
			Notification: "cannot be empty",
		}

		emptyPassword = true
		userErrors = append(userErrors, userError)
	}

	if !emptyPassword {
		if !PasswordsMatch(userCreate.Password, userCreate.PasswordConfirm) {
			userError := &domainError.DomainValidationError{
				Field:        "password",
				FieldType:    "required",
				Notification: "passwords do not match",
			}

			passwordsMatch = false
			userErrors = append(userErrors, userError)
		}
	}

	if passwordsMatch {
		if domainValidatorUtility.IsStringLengthExceeded(userCreate.Password) {
			userError := &domainError.DomainValidationError{
				Field:        "password",
				FieldType:    "required",
				Notification: "cannot be more than 100 characters long",
			}

			userErrors = append(userErrors, userError)
		}

		if IsPasswordStringContainsSpecialCharacters(userCreate.Password) {
			userError := &domainError.DomainValidationError{
				Field:        "password",
				FieldType:    "required",
				Notification: "sorry, only letters (a-z), numbers(0-9), the hyphen and the underscore characters are allowed",
			}

			userErrors = append(userErrors, userError)
		}
	}

	return userErrors
}

func (userUpdate *UserUpdate) UserUpdateValidator() error {
	var message string
	var err error

	if userUpdate.Name == "" {
		message = "key: `UserUpdateView.Name` error: field validation for `name` failed, `name` cannot be empty "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userLogin *UserLogin) UserLoginValidator() error {
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

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userForgottenPassword *UserForgottenPassword) UserForgottenPasswordValidator() error {
	var message string
	var err error

	if userForgottenPassword.Email == "" {
		message = "key: `UserForgottenPasswordView.Email` error: field validation for `email` failed, `email` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userForgottenPassword.Email != "" {
		_, ok := mail.ParseAddress(userForgottenPassword.Email)
		if ok != nil {
			message = message + "key: `UserForgottenPasswordView.Email` error: field validation for `email` failed, invalid email address "
			err = fmt.Errorf(message)
		}
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func (userResetPassword *UserResetPassword) UserResetPasswordValidator() error {
	var message string
	var err error

	if userResetPassword.Password == "" {
		message = "key: `UserResetPasswordView.Password` error: field validation for `password` failed, `password` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userResetPassword.PasswordConfirm == "" {
		message = "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, `password_confirm` cannot be empty "
		err = fmt.Errorf(message)
	}

	if userResetPassword.Password != "" && userResetPassword.PasswordConfirm != "" && userResetPassword.Password != userResetPassword.PasswordConfirm {
		message = message + "key: `UserResetPasswordView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match "
		err = fmt.Errorf(message)
	}

	if err != nil {
		message = strings.TrimSpace(message)
		err = fmt.Errorf(message)

		return err
	}

	return nil
}

func UserValidator(checkedStringKey string, checkedStringValue string, jsonStringKey string) string {
	var message string

	if checkedStringValue == "" {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" + " cannot be empty "
	}

	if len(checkedStringValue) > 100 {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" + " cannot be more than 100 characters long "
	}

	if !regexp.MustCompile(`^[a-zA-z0-9- \t]*$`).MatchString(checkedStringValue) {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" can use only letters, numbers, spaces, the hyphen or underscore character "
	}

	return message
}

func UserPasswordValidator(checkedStringKey string, checkedStringValue string, jsonStringKey string) string {
	var message string

	if checkedStringValue == "" {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" cannot be empty "
	}

	if len(checkedStringValue) > 100 {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" cannot be more than 100 characters long "
	}

	if !regexp.MustCompile(`^[a-zA-z0-9-]*$`).MatchString(checkedStringValue) {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" can use only letters, numbers, the hyphen or underscore character "
	}

	return message
}

func IsUserNameContainsSpecialCharacters(checkedString string) bool {
	flag := false

	if !regexp.MustCompile(`^[a-zA-z0-9-_ \t]*$`).MatchString(checkedString) {
		flag = true
	}

	return flag
}

func IsEmailValid(emailString string) bool {
	flag := true

	if !regexp.MustCompile(emailRegexString).MatchString(emailString) {
		flag = false
	}

	return flag
}

func IsEmailDomainValid(emailString string) bool {
	flag := true

	if len(emailString) == 0 {
		return true
	}

	host := strings.Split(emailString, "@")[1]

	_, err := net.LookupMX(host)

	if err != nil {
		flag = false
	}

	return flag
}

func PasswordsMatch(password string, passwordConfirm string) bool {
	flag := true

	if password != passwordConfirm {
		flag = false
	}

	return flag
}

func IsPasswordStringContainsSpecialCharacters(checkedString string) bool {
	flag := false

	if !regexp.MustCompile(`^[a-zA-z0-9-_,.]*$`).MatchString(checkedString) {
		flag = true
	}

	return flag
}
