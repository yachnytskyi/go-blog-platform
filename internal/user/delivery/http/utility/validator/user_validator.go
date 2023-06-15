package validator

import (
	"net"
	"regexp"
	"strings"
)

const (
	EmailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[\\\\\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[\\\\+\\-\\/=\\\\_{\\|}]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.||[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

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

func UserEmailValidator(checkedStringKey string, checkedStringValue string, jsonStringKey string) string {
	var message string

	if checkedStringValue == "" {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" +
			" failed, " + "'" + jsonStringKey + "'" + " cannot be empty "
	}

	if len(checkedStringValue) > 100 {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" cannot be more than 100 characters long "
	}

	if !regexp.MustCompile(EmailRegexString).MatchString(checkedStringValue) {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" is invalid "
	}

	host := strings.Split(checkedStringValue, "@")[1]

	_, err := net.LookupMX(host)
	if err != nil {
		message = "key: " + "`" + checkedStringKey + "`" + " error: field validation for " + "`" + jsonStringKey + "`" + " failed, " + "'" + jsonStringKey + "'" +
			" domain is not found "
	}

	return message
}

func UserPasswordMatchValidator(password string, passwordConfirm string) string {
	var message string

	if password != passwordConfirm {
		message = "key: `UserCreateView.PasswordConfirm` error: field validation for `password_confirm` failed, passwords do not match "

		return message
	}

	return message
}
