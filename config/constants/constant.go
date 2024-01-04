package constants

import (
	"net/http"
	"time"
)

const (
	// Context timers.
	DefaultContextTimer = time.Duration(time.Second * 5)

	// Tokens
	AccessTokenValue  = "access_token"
	RefreshTokenValue = "refresh_token"
	LoggedInValue     = "logged_in"
	TokenDomainValue  = "localhost"
	LogoutMaxAgeValue = -1
	UserIDContext     = "userID"
	UserContext       = "user"

	// Email subjects and URLs.
	EmailConfirmationUrl     string = "users/verifyemail/"
	ForgottenPasswordUrl     string = "users/reset-password/"
	EmailConfirmationSubject string = "Your account verification code"
	ForgottenPasswordSubject string = "Your password reset token (it is valid for 15 minutes)"

	// Pagination.
	DefaultPage      = "1"
	DefaultLimit     = "10"
	MaxItemsPerPage  = 100
	DefaultOrderBy   = "created_at"
	DefaultSortOrder = "descend"
	SortAscend       = "ascend"
	SortDescend      = "descent"

	// Regex patterns.
	StringRegex             = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`
	TitleStringRegex        = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`
	TextStringRegex         = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$`
	MinStringLength         = 4
	MaxStringLength         = 40
	MinOptionalStringLength = 0
	FieldRequired           = "required"
	FieldOptional           = "optional"
	EmptyString             = ""
	True                    = "true"
	False                   = "false"

	// User Notifications.
	SendingEmailNotification                 = "We have sent an email with a verification code to the provided address "
	SendingEmailWithInstructionsNotification = "You will receive an email with detailed instructions shortly."

	// HTTP codes.
	StatusOk           = http.StatusOK
	StatusCreated      = http.StatusCreated
	StatusNoContent    = http.StatusNoContent
	StatusBadRequest   = http.StatusBadRequest
	StatusUnauthorized = http.StatusUnauthorized
	StatusForbidden    = http.StatusForbidden
	StatusBadGateway   = http.StatusBadGateway

	// Error Messages.
	StringAllowedLength             = "Can be between %d and %d characters long."
	StringOptionalAllowedLength     = "Cannot be more than %d characters."
	EmailAlreadyExists              = "An account with this email address already exists."
	AuthorizationErrorNotification  = "Access denied. You do not have the required permissions to perform this action. Please try again or contact our support team for assistance."
	LoggingErrorNotification        = "You are not logged in."
	AlreadyRegisteredNotification   = "You are already registered, and registration is not allowed for existing users."
	EntityNotFoundErrorNotification = "Sorry, the requested item does not exist in our records."
	PaginationErrorNotification     = "Sorry, there was an issue with the pagination request. Please check your parameters and try again."
	InternalErrorNotification       = "Oops! Something went wrong on our end. Please try again later or contact our support team for assistance."
	InvalidHTTPMethodNotification   = "Invalid HTTP method"
	InvalidContentTypeNotification  = "Invalid content type"
	InvalidHeaderFormat             = "Invalid header format"

	// Databases.
	MongoDB = "MongoDB"

	// Domains.
	UseCase = "UseCase"

	// Deliveries
	Gin = "Gin"
)
