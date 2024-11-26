package usecase

import (
	"fmt"
	"net"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	bcrypt "golang.org/x/crypto/bcrypt"
)

// Constants used for various validation messages and field names.
const (
	// Regex Patterns for validating email, username, and password.
	emailRegex    = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\\\/=\\\\{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	usernameRegex = `^[a-zA-z0-9-_ \t]*$`
	passwordRegex = `^[a-zA-z0-9-_*,.]*$`

	// Error Messages for invalid inputs.
	passwordAllowedCharacters = "Sorry, only letters (a-z), numbers(0-9), the asterics, hyphen and underscore characters are allowed."
	emailAllowedCharacters    = "Sorry, only letters (a-z), numbers(0-9) and periods (.) are allowed, you cannot use a period in the end and more than one in a row."
	invalidEmailDomain        = "Email domain does not exist."
	passwordsDoNotMatch       = "Passwords do not match."
	invalidEmailOrPassword    = "Invalid email or password."

	// Field Names used in validation.
	usernameField         = "username"
	EmailField            = "email"
	passwordField         = "password"
	emailOrPasswordFields = "email or password"
	resetTokenField       = "reset token"
)

var (
	emailValidator = utility.StringValidator{
		FieldName:    EmailField,
		FieldRegex:   emailRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: emailAllowedCharacters,
		IsOptional:   false,
	}
	usernameValidator = utility.StringValidator{
		FieldName:    usernameField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: constants.StringAllowedCharacters,
		IsOptional:   false,
	}
	passwordValidator = utility.StringValidator{
		FieldName:    passwordField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: passwordAllowedCharacters,
		IsOptional:   false,
	}
	tokenValidator = utility.StringValidator{
		FieldName:  resetTokenField,
		FieldRegex: usernameRegex,
		MinLength:  resetTokenLength,
		MaxLength:  resetTokenLength,
		IsOptional: false,
	}
	// Add more validators for other fields as needed.
)

func validateUserCreate(logger interfaces.Logger, userCreate user.UserCreate) common.Result[user.UserCreate] {
	validationErrors := make([]error, 0, 4)
	userCreate.Email = commonUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Username = strings.TrimSpace(userCreate.Username)
	userCreate.Password = strings.TrimSpace(userCreate.Password)
	userCreate.PasswordConfirm = strings.TrimSpace(userCreate.PasswordConfirm)
	emailValidator.Field = userCreate.Email
	usernameValidator.Field = userCreate.Username
	
	validationErrors = validateEmail(logger, location+"validateUserCreate", userCreate.Email, validationErrors)
	validationErrors = utility.ValidateField(logger, location+"validateUserCreate", usernameValidator, validationErrors)
	validationErrors = validatePassword(logger, location+"validateUserCreate", userCreate.Password, userCreate.PasswordConfirm, validationErrors)
	if len(validationErrors) > 0 {
		return common.NewResultOnFailure[user.UserCreate](domain.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[user.UserCreate](userCreate)
}

func validateUserUpdate(logger interfaces.Logger, userUpdate user.UserUpdate) common.Result[user.UserUpdate] {
	validationErrors := make([]error, 0, 1)
	userUpdate.Username = strings.TrimSpace(userUpdate.Username)
	usernameValidator.Field = userUpdate.Username

	validationErrors = utility.ValidateField(logger, location+"validateUserUpdate", usernameValidator, validationErrors)
	if len(validationErrors) > 0 {
		return common.NewResultOnFailure[user.UserUpdate](domain.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[user.UserUpdate](userUpdate)
}

func validateUserLogin(logger interfaces.Logger, userLogin user.UserLogin) common.Result[user.UserLogin] {
	validationErrors := make([]error, 0, 2)
	userLogin.Email = commonUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = strings.TrimSpace(userLogin.Password)
	passwordValidator.Field = userLogin.Password

	validationErrors = validateEmail(logger, location+"validateUserLogin", userLogin.Email, validationErrors)
	validationErrors = utility.ValidateField(logger, location+"validateUserLogin", passwordValidator, validationErrors)
	if len(validationErrors) > 0 {
		return common.NewResultOnFailure[user.UserLogin](domain.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[user.UserLogin](userLogin)
}

func validateUserForgottenPassword(logger interfaces.Logger, userForgottenPassword user.UserForgottenPassword) common.Result[user.UserForgottenPassword] {
	validationErrors := make([]error, 0, 2)
	userForgottenPassword.Email = commonUtility.SanitizeAndToLowerString(userForgottenPassword.Email)

	validationErrors = validateEmail(logger, location+"validateUserForgottenPassword", userForgottenPassword.Email, validationErrors)
	if len(validationErrors) > 0 {
		return common.NewResultOnFailure[user.UserForgottenPassword](domain.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[user.UserForgottenPassword](userForgottenPassword)
}

func validateUserResetPassword(logger interfaces.Logger, userResetPassword user.UserResetPassword) common.Result[user.UserResetPassword] {
	validationErrors := make([]error, 0, 2)
	userResetPassword.ResetToken =  strings.TrimSpace(userResetPassword.ResetToken)
	tokenValidator.Field = userResetPassword.ResetToken
	userResetPassword.Password = strings.TrimSpace(userResetPassword.Password)
	userResetPassword.PasswordConfirm = strings.TrimSpace(userResetPassword.PasswordConfirm)

	validationErrors = validatePassword(logger, location+"validateUserResetPassword", userResetPassword.Password, userResetPassword.PasswordConfirm, validationErrors)
	validationErrors = utility.ValidateField(logger, location+"validateUserResetPassword", tokenValidator, validationErrors)
	if len(validationErrors) > 0 {
		return common.NewResultOnFailure[user.UserResetPassword](domain.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[user.UserResetPassword](userResetPassword)
}

func validateEmail(logger interfaces.Logger, location, email string, validationErrors []error) []error {
	errors := validationErrors
	emailValidator.Field = email

	validateFieldError := validateField(logger, location+".validateEmail", emailValidator)
	if validator.IsError(validateFieldError) {
		errors = append(errors, validateFieldError)
		return errors
	}

	checkEmailError := checkEmailDomain(logger, location+".validateEmail", email)
	if validator.IsError(checkEmailError) {
		errors = append(errors, checkEmailError)
	}

	return errors
}

func validatePassword(logger interfaces.Logger, location, password, passwordConfirm string, validationErrors []error) []error {
	errors := validationErrors
	passwordValidator.Field = password

	validateFieldError := validateField(logger, location+".validatePassword", passwordValidator)
	if validator.IsError(validateFieldError) {
		errors = append(errors, validateFieldError)
		return errors
	}
	if password != passwordConfirm {
		validationError := domain.NewValidationError(
			location+".validatePassword",
			passwordValidator.FieldName,
			constants.FieldRequired,
			passwordsDoNotMatch,
		)

		logger.Debug(validationError)
		errors = append(errors, validationError)
	}

	return errors
}

// checkEmailDomain checks if the email domain exists by resolving DNS records.
func checkEmailDomain(logger interfaces.Logger, location, emailString string) error {
	host := strings.Split(emailString, "@")[1]
	_, lookupMXError := net.LookupMX(host)
	if validator.IsError(lookupMXError) {
		validationError := domain.NewValidationError(
			location+".checkEmailDomain",
			EmailField,
			constants.FieldRequired,
			invalidEmailDomain,
		)

		logger.Debug(validationError)
		return validationError
	}

	return nil
}

func checkPasswords(logger interfaces.Logger, location, hashedPassword string, checkedPassword string) error {
	if validator.IsError(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		validationError := domain.NewValidationError(
			location+".checkPasswords.CompareHashAndPassword",
			emailOrPasswordFields,
			constants.FieldRequired,
			passwordsDoNotMatch,
		)
		logger.Debug(validationError)
		validationError.Notification = invalidEmailOrPassword
		return validationError
	}

	return nil
}

func checkEmail(logger interfaces.Logger, location, email string) error {
	emailValidator.Field = email
	validateFieldError := validateField(logger, location+".checkEmail", emailValidator)
	if validator.IsError(validateFieldError) {
		return validateFieldError
	}

	return checkEmailDomain(logger, location+".checkEmail", email)
}

func validateField(logger interfaces.Logger, location string, stringValidator utility.StringValidator) error {
	if utility.IsStringLengthInvalid(stringValidator.Field, stringValidator.MinLength, stringValidator.MaxLength) {
		validationError := domain.NewValidationError(
			location+".validateField.IsStringLengthInvalid",
			stringValidator.FieldName,
			constants.FieldRequired,
			fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength),
		)
		logger.Debug(validationError)
		return validationError
	}
	if utility.AreStringCharactersInvalid(stringValidator.Field, stringValidator.FieldRegex) {
		validationError := domain.NewValidationError(
			location+".validateField.AreStringCharactersInvalid",
			stringValidator.FieldName,
			constants.FieldRequired,
			stringValidator.Notification,
		)

		logger.Debug(validationError)
		return validationError
	}

	return nil
}
