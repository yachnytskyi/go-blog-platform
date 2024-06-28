package constants

import (
	"net/http"
	"time"
)

// Schemes used in the application.
const (
	HTTP  = "http"  // HTTP scheme.
	HTTPS = "https" // HTTPS scheme.
)

// HTTP Headers used in the application.
const (
	ContentType   = "Content-Type"  // Content-Type header.
	Authorization = "Authorization" // Authorization header.
	Bearer        = "Bearer"        // Bearer token prefix.
)

// Pagination constants.
const (
	PageValue            = "page"       // Page query parameter.
	LimitValue           = "limit"      // Limit query parameter.
	OrderByValue         = "order_by"   // Order by query parameter.
	SortOrderValue       = "sort_order" // Sort order query parameter.
	DefaultAmountOfPages = 10           // Default amount of pages for pagination.
	DefaultPage          = "1"          // Default page number.
	DefaultLimit         = "10"         // Default limit per page.
	MaxItemsPerPage      = 100          // Maximum items per page.
	DefaultOrderBy       = "created_at" // Default order by field.
	DefaultSortOrder     = "descend"    // Default sort order.
	SortAscend           = "ascend"     // Ascending sort order.
	SortDescend          = "descend"    // Descending sort order.
)

// Context timers.
const (
	DefaultContextTimer = 5 * time.Second // Default context timeout duration.
)

// DateTime format.
const (
	DateTimeFormat = "02-Jan-2006 03:04 PM MST" // DateTime format string.
)

// Token-related constants.
const (
	AccessTokenValue             = "access_token"                          // Access token value.
	RefreshTokenValue            = "refresh_token"                         // Refresh token value.
	LoggedInValue                = "logged_in"                             // Logged in status value.
	TokenDomainValue             = "localhost"                             // Token domain value.
	LogoutMaxAgeValue            = -1                                      // Logout max age value.
	User              contextKey = "user"                                  // User context key.
	ID                contextKey = "id"                                    // ID context key.
	UserRole          contextKey = "userRole"                              // User role context key.
	IDContextMissing             = "ID context value is missing or empty." // ID context missing error message.
)

// Email subjects and URLs.
const (
	EmailConfirmationUrl     = "users/verifyemail/"                                     // Email confirmation URL.
	ForgottenPasswordUrl     = "users/reset-password/"                                  // Forgotten password URL.
	EmailConfirmationSubject = "Your account verification code"                         // Email confirmation subject.
	ForgottenPasswordSubject = "Your password reset token (it is valid for 15 minutes)" // Forgotten password subject.
)

// Regex patterns and string length constraints.
const (
	StringRegex             = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`     // General string regex pattern.
	TitleStringRegex        = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`                  // Title string regex pattern.
	TextStringRegex         = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$` // Text string regex pattern.
	MinStringLength         = 4                                                      // Minimum string length.
	MaxStringLength         = 40                                                     // Maximum string length.
	MinOptionalStringLength = 0                                                      // Minimum optional string length.
	MaxOptionalStringLength = 40                                                     // Maximum optional string length.
	FieldRequired           = "required"                                             // Field required status.
	FieldOptional           = "optional"                                             // Field optional status.
	True                    = "true"                                                 // True value.
	False                   = "false"                                                // False value.
)

// User Notifications.
const (
	SendingEmailNotification                 = "We have sent an email with a verification code to the provided address."
	SendingEmailWithInstructionsNotification = "You will receive an email with detailed instructions shortly."
)

// HTTP status codes.
const (
	StatusOk               = http.StatusOK               // HTTP 200 OK.
	StatusCreated          = http.StatusCreated          // HTTP 201 Created.
	StatusNoContent        = http.StatusNoContent        // HTTP 204 No Content.
	StatusBadRequest       = http.StatusBadRequest       // HTTP 400 Bad Request.
	StatusUnauthorized     = http.StatusUnauthorized     // HTTP 401 Unauthorized.
	StatusForbidden        = http.StatusForbidden        // HTTP 403 Forbidden.
	StatusNotFound         = http.StatusNotFound         // HTTP 404 Not Found.
	StatusMethodNotAllowed = http.StatusMethodNotAllowed // HTTP 405 Not Allowed.
	StatusBadGateway       = http.StatusBadGateway       // HTTP 502 Bad Gateway.
)

// HTTP methods for RESTful operations.
const (
	Get    = "GET"    // GET method for retrieving resources.
	Post   = "POST"   // POST method for creating resources.
	Put    = "PUT"    // PUT method for updating resources.
	Patch  = "PATCH"  // PATCH method for partially updating resources.
	Delete = "DELETE" // DELETE method for deleting resources.
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
	PaginationErrorNotification    = "Sorry, there was an issue with the pagination request. Please check your parameters and try again."                                           // Pagination error message.
	InternalErrorNotification      = "Oops! Something went wrong on our end. Please try again later or contact our support team for assistance."                                    // Internal error message.
	InvalidHTTPMethodNotification  = "Invalid HTTP method. You can only use the methods from the following list: "                                                                  // Invalid HTTP method message.
	InvalidContentTypeNotification = "Invalid content type. You can use only them from the following list: "                                                                        // Invalid content type message.
	InvalidHeaderFormat            = "Invalid header format."                                                                                                                       // Invalid header format message.
)

// Databases used in the application.
const (
	MongoDB = "MongoDB" // MongoDB database name.
)

// Domains used in the application.
const (
	UseCase = "UseCaseV1" // UseCase domain name.
)

// Deliveries used in the application.
const (
	Gin = "Gin" // Gin delivery name.
)

// Common routes used in the application.
const (
	GetAllItemsURL = ""     // Endpoint for fetching all items.
	GetItemByIdURL = "/:id" // Endpoint for fetching an item by its ID.
	ItemIdParam    = "id"   // Parameter name for item ID.
)

// Domain-specific routes.
const (
	UsersGroupPath = "/users" // Users domain route.
	// Initialize other routes here.
)

// Database table names.
const (
	UsersTable = "users" // Users table name in the database.
)

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string
