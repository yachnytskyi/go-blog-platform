package usecase

import (
	"fmt"
	"net"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	domainValidator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator/domain"
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
	usernameField         = "name"
	EmailField            = "email"
	passwordField         = "password"
	emailOrPasswordFields = "email or password"
)

// Validators for email, username, and password fields.
var (
	emailValidator = domainModel.CommonValidator{
		FieldName:  EmailField,
		FieldRegex: emailRegex,
		MinLength:  constants.MinStringLength,
		MaxLength:  constants.MaxStringLength,
	}
	usernameValidator = domainModel.CommonValidator{
		FieldName:    usernameField,
		FieldRegex:   usernameRegex,
		MinLength:    constants.MinStringLength,
		MaxLength:    constants.MaxStringLength,
		Notification: constants.StringAllowedCharacters,
	}
	passwordValidator = domainModel.CommonValidator{
		FieldName:  passwordField,
		FieldRegex: usernameRegex,
		MinLength:  constants.MinStringLength,
		MaxLength:  constants.MaxStringLength,
	}
	// Add more validators for other fields as needed.
)

// validateUserCreate validates the fields of the UserCreate struct.
func validateUserCreate(userCreate userModel.UserCreate) common.Result[userModel.UserCreate] {
	// Initialize a slice to hold validation errors.
	validationErrors := make([]error, 0, 4)

	// Sanitize input fields.
	userCreate.Email = domainUtility.SanitizeAndToLowerString(userCreate.Email)
	userCreate.Name = domainUtility.SanitizeString(userCreate.Name)
	userCreate.Password = domainUtility.SanitizeString(userCreate.Password)
	userCreate.PasswordConfirm = domainUtility.SanitizeString(userCreate.PasswordConfirm)

	// Perform validation for each field.
	validationErrors = validateEmail(userCreate.Email, validationErrors)
	validationErrors = domainValidator.ValidateField(userCreate.Name, usernameValidator, validationErrors)
	validationErrors = validatePassword(userCreate.Password, userCreate.PasswordConfirm, validationErrors)

	// Return validation result based on presence of errors.
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserCreate](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserCreate](userCreate)
}

// validateUserUpdate validates the fields of the UserUpdate struct.
func validateUserUpdate(userUpdate userModel.UserUpdate) common.Result[userModel.UserUpdate] {
	// Initialize a slice to hold validation errors.
	validationErrors := make([]error, 0, 1)

	// Sanitize input fields.
	userUpdate.Name = domainUtility.SanitizeString(userUpdate.Name)

	// Perform validation for each field.
	validationErrors = domainValidator.ValidateField(userUpdate.Name, usernameValidator, validationErrors)

	// Return validation result based on presence of errors.
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserUpdate](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserUpdate](userUpdate)
}

// validateUserLogin validates the fields of the UserLogin struct.
func validateUserLogin(userLogin userModel.UserLogin) common.Result[userModel.UserLogin] {
	// Initialize a slice to hold validation errors.
	validationErrors := make([]error, 0, 2)

	// Sanitize input fields.
	userLogin.Email = domainUtility.SanitizeAndToLowerString(userLogin.Email)
	userLogin.Password = domainUtility.SanitizeString(userLogin.Password)

	// Perform validation for each field.
	validationErrors = validateEmail(userLogin.Email, validationErrors)
	validationErrors = domainValidator.ValidateField(userLogin.Password, usernameValidator, validationErrors)

	// Return validation result based on presence of errors.
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserLogin](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserLogin](userLogin)
}

// validateUserForgottenPassword validates the fields of the UserForgottenPassword struct.
func validateUserForgottenPassword(userForgottenPassword userModel.UserForgottenPassword) common.Result[userModel.UserForgottenPassword] {
	// Initialize a slice to hold validation errors.
	validationErrors := make([]error, 0, 1)

	// Sanitize input fields.
	userForgottenPassword.Email = domainUtility.SanitizeAndToLowerString(userForgottenPassword.Email)

	// Perform validation for the email field.
	validationErrors = validateEmail(userForgottenPassword.Email, validationErrors)

	// Return validation result based on presence of errors.
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserForgottenPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserForgottenPassword](userForgottenPassword)
}

// validateResetPassword validates the fields of the UserResetPassword struct.
func validateResetPassword(userResetPassword userModel.UserResetPassword) common.Result[userModel.UserResetPassword] {
	// Initialize a slice to hold validation errors.
	validationErrors := make([]error, 0, 2)

	// Sanitize input fields.
	userResetPassword.Password = domainUtility.SanitizeString(userResetPassword.Password)
	userResetPassword.PasswordConfirm = domainUtility.SanitizeString(userResetPassword.PasswordConfirm)

	// Perform validation for each field.
	validationErrors = validatePassword(userResetPassword.Password, userResetPassword.PasswordConfirm, validationErrors)

	// Return validation result based on presence of errors.
	if validator.IsSliceNotEmpty(validationErrors) {
		return common.NewResultOnFailure[userModel.UserResetPassword](domainError.NewValidationErrors(validationErrors))
	}

	return common.NewResultOnSuccess[userModel.UserResetPassword](userResetPassword)
}

// validateEmail validates the email field using regex and other checks.
func validateEmail(email string, validationErrors []error) []error {
	// Initialize the errors slice with the given validationErrors to accumulate any validation errors.
	// This allows appending new errors while preserving the existing ones.
	errors := validationErrors

	// Check email length.
	if domainValidator.IsStringLengthInvalid(email, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthInvalid", EmailField, constants.FieldRequired, notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check email characters.
	if domainValidator.AreStringCharactersInvalid(email, emailValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"validateEmail.AreStringCharactersInvalid", EmailField, constants.FieldRequired, emailAllowedCharacters)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check email domain validity.
	if isEmailDomainNotValid(email) {
		validationError := domainError.NewValidationError(location+"validateEmail.IsEmailDomainNotValid", EmailField, constants.FieldRequired, invalidEmailDomain)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

// validatePassword validates the password and password confirmation fields.
func validatePassword(password, passwordConfirm string, validationErrors []error) []error {
	// Initialize the errors slice with the given validationErrors to accumulate any validation errors.
	// This allows appending new errors while preserving the existing ones.
	errors := validationErrors

	// Check password length.
	if domainValidator.IsStringLengthInvalid(password, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, passwordValidator.MinLength, passwordValidator.MaxLength)
		validationError := domainError.NewValidationError(location+"validatePassword.IsStringLengthInvalid", passwordValidator.FieldName, constants.FieldRequired, notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check password characters.
	if domainValidator.AreStringCharactersInvalid(password, passwordValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringCharactersInvalid", passwordValidator.FieldName, constants.FieldRequired, passwordAllowedCharacters)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check if passwords match.
	if validator.AreStringsNotEqual(password, passwordConfirm) {
		validationError := domainError.NewValidationError(location+"validatePassword.AreStringsNotEqual", passwordValidator.FieldName, constants.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

// isEmailDomainNotValid checks if the email domain is valid by performing an MX record lookup.
func isEmailDomainNotValid(emailString string) bool {
	host := strings.Split(emailString, "@")[1]
	_, lookupMXError := net.LookupMX(host)
	return validator.IsError(lookupMXError)
}

// arePasswordsNotEqual compares the hashed password with the provided password.
func arePasswordsNotEqual(hashedPassword string, checkedPassword string) error {
	if validator.IsError(bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(checkedPassword))) {
		validationError := domainError.NewValidationError(location+"arePasswordsNotEqual.CompareHashAndPassword", emailOrPasswordFields, constants.FieldRequired, passwordsDoNotMatch)
		logging.Logger(validationError)
		validationError.Notification = invalidEmailOrPassword
		return validationError
	}

	return nil
}

// isEmailValid performs basic email validation checks.
func isEmailValid(email string) error {
	if domainValidator.IsStringLengthInvalid(email, constants.MinStringLength, constants.MaxStringLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, constants.MinStringLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"validateEmail.IsStringLengthInvalid", EmailField, constants.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}

	if domainValidator.AreStringCharactersInvalid(email, emailValidator.FieldRegex) {
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
