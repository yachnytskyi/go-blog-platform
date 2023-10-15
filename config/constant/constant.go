package constant

import "time"

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

	// Pagination.
	DefaultPage     = "1"
	DefaultLimit    = "10"
	MaxItemsPerPage = 100

	// Regex patterns.
	StringRegex             = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`
	TitleStringRegex        = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`
	TextStringRegex         = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$`
	MinStringLength         = 4
	MaxStringLength         = 40
	MinOptionalStringLength = 0
	FieldRequired           = "required"
	FieldOptional           = "optional"

	// User Notifications.
	SendingEmailNotification                   = "We sent an email with a verification code to "
	SendingEmailWithIntstructionsNotifications = "We sent you an email with needed instructions."

	// Error Messages.
	StringAllowedLength             = "can be between %d and %d characters long."
	StringOptionalAllowedLength     = "cannot be more than %d characters."
	EmailAlreadyExists              = "user with this email already exists."
	InternalErrorNotification       = "something went wrong, please repeat later."
	EntityNotFoundErrorNotification = "this item does not exist."

	// Databases.
	MongoDB = "MongoDB"

	// Domains.
	UseCase = "UseCase"

	// Deliveries
	Gin = "Gin"
)
