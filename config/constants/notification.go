package constants

// Unsupported Types.
const (
	UnsupportedConfig     = "Unsupported config type: %s"     // Unsupported config type error message.
	UnsupportedLogger     = "Unsupported logger type: %s"     // Unsupported logger type error message.
	UnsupportedRepository = "Unsupported repository type: %s" // Unsupported repository type error message.
	UnsupportedUseCase    = "Unsupported use case type: %s"   // Unsupported use case type error message.
	UnsupportedDelivery   = "Unsupported delivery type: %s"   // Unsupported delivery type error message.
)

// Server Notifications.
const (
	DatabaseConnectionSuccess = "Database connection is established..."               // Database connection success message.
	DatabaseConnectionClosed  = "Database connection has been successfully closed..." // Database connection closed message.
	DatabaseConnectionFailure = "Failed to establish database connection"             // Database connection failure message.
	ServerConnectionSuccess   = "Server is successfully launched..."                  // Server connection success message.
	ServerConnectionClosed    = "Server has been successfully shut down..."           // Server connection closed message.
)

// User Notifications.
const (
	LogoutNotificationMessage                = "You are successfully logged out."                                               // Logout success message.
	SendingEmailNotification                 = "We have sent an email with a verification code to the provided address."        // Email sent notification.
	SendingEmailWithInstructionsNotification = "You will receive an email with detailed instructions shortly."                  // Email with instructions sent notification.
	PasswordResetSuccessNotification         = "Congratulations! Your password was updated successfully! Please sign in again." // Password reset success message.
)

// Error Messages.
const (
	StringAllowedLength            = "Can be between %d and %d characters long."                                                                                                    // Allowed string length message.
	StringOptionalAllowedLength    = "Can be empty or between %d and %d characters long."                                                                                           // Optional string length message.
	StringAllowedCharacters        = "Sorry, only letters (a-z), numbers (0-9), and spaces are allowed."                                                                            // Allowed string character message.
	EmailAlreadyExists             = "An account with this email address already exists."                                                                                           // Email already exists message.
	EmailTemplateNotFound          = "Email template not found."                                                                                                                    // Email template not found message.
	AuthorizationErrorNotification = "Access denied. You do not have the required permissions to perform this action. Please try again or contact our support team for assistance." // Authorization error message.
	LoggingErrorNotification       = "You are not logged in."                                                                                                                       // Not logged in message.
	MethodNotAllowedNotification   = "Method %s is not allowed."                                                                                                                    // Method not allowed message.
	RouteNotFoundNotification      = "The requested URL '%s' was not found on this server."                                                                                         // Route not found message.
	AlreadyLoggedInNotification    = "Already logged in. This action is not allowed."                                                                                               // Already logged in message.
	ItemNotFoundErrorNotification  = "Sorry, the requested item does not exist in our records."                                                                                     // Item not found message.
	TimeExpiredErrorNotification   = "Sorry, the time is expired and not valid anymore"                                                                                             // Time expired message.
	PaginationErrorNotification    = "Sorry, there was an issue with the pagination request. Please check your parameters and try again."                                           // Pagination error message.
	InternalErrorNotification      = "Oops! Something went wrong on our end. Please try again later or contact our support team for assistance."                                    // Internal error message.
	InvalidHTTPMethodNotification  = "Invalid HTTP method. You can only use the methods from the following list: "                                                                  // Invalid HTTP method message.
	InvalidContentTypeNotification = "Invalid content type. You can use only them from the following list: "                                                                        // Invalid content type message.
	InvalidHeaderFormat            = "Invalid header format."                                                                                                                       // Invalid header format message.
	InvalidTokenErrorMessage       = "The token is invalid. Please use the correct token."                                                                                          // Error message for invalid tokens.
)
