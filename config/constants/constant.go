package constants

import (
	"net/http"
	"time"
)

// Schemes.
const (
	HTTP  = "http"  // HTTP scheme.
	HTTPS = "https" // HTTPS scheme.
)

// Headers.
const (
	ContentType = "Content-Type" // Content-Type header.
)

// Pagination
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

// Context timers
const (
	DefaultContextTimer = 5 * time.Second // Default context timeout duration.
)

// Datetime format
const (
	DateTimeFormat = "02-Jan-2006 03:04 PM MST" // DateTime format string.
)

// Tokens.
const (
	AccessTokenValue             = "access_token"  // Access token value.
	RefreshTokenValue            = "refresh_token" // Refresh token value.
	LoggedInValue                = "logged_in"     // Logged in status value.
	TokenDomainValue             = "localhost"     // Token domain value.
	LogoutMaxAgeValue            = -1              // Logout max age value.
	UserContext       contextKey = "user"          // User  context key.
	UserIDContext     contextKey = "userID"        // User ID context key.
	UserRoleContext   contextKey = "userRole"      // User role context key.
)

// Email subjects and URLs.
const (
	EmailConfirmationUrl     = "users/verifyemail/"                                     // Email confirmation URL.
	ForgottenPasswordUrl     = "users/reset-password/"                                  // Forgotten password URL.
	EmailConfirmationSubject = "Your account verification code"                         // Email confirmation subject.
	ForgottenPasswordSubject = "Your password reset token (it is valid for 15 minutes)" // Forgotten password subject.
)

// Regex patterns.
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
	SendingEmailNotification                 = "We have sent an email with a verification code to the provided address "
	SendingEmailWithInstructionsNotification = "You will receive an email with detailed instructions shortly."
)

// HTTP codes.
const (
	StatusOk               = http.StatusOK               // HTTP 200 OK.
	StatusCreated          = http.StatusCreated          // HTTP 201 Created.
	StatusNoContent        = http.StatusNoContent        // HTTP 204 No Content.
	StatusBadRequest       = http.StatusBadRequest       // HTTP 400 Bad Request.
	StatusUnauthorized     = http.StatusUnauthorized     // HTTP 401 Unauthorized.
	StatusForbidden        = http.StatusForbidden        // HTTP 403 Forbidden.
	StatusNotFound         = http.StatusNotFound         // HTTP 404 Not Found.
	StatusMethodNotAllowed = http.StatusMethodNotAllowed /// HTTP 405 Not Allowed.
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
	StringAllowedLength            = "Can be between %d and %d characters long."          // Allowed string length message.
	StringOptionalAllowedLength    = "Cannot be more than %d characters."                 // Optional string length message.
	EmailAlreadyExists             = "An account with this email address already exists." // Email already exists message.
	EmailTemplateNotFound          = "Email template not found"                           // Email template not found message.
	AuthorizationErrorNotification = "Access denied. You do not have the required permissions to perform this action. Please try again or contact our support team for assistance."
	LoggingErrorNotification       = "You are not logged in."                                                                                    // Not logged in message.
	MethodNotAllowedNotification   = "Method %s is not allowed."                                                                                 // Method not allowed in message.
	RouteNotFoundNotification      = "The requested URL '%s' was not found on this server."                                                      // Route not found message.
	AlreadyLoggedInNotification    = "Already logged in. This action is not allowed."                                                            // Already logged in message.
	ItemNotFoundErrorNotification  = "Sorry, the requested item does not exist in our records."                                                  // Item not found message.
	PaginationErrorNotification    = "Sorry, there was an issue with the pagination request. Please check your parameters and try again."        // Pagination error message.
	InternalErrorNotification      = "Oops! Something went wrong on our end. Please try again later or contact our support team for assistance." // Internal error message.
	InvalidHTTPMethodNotification  = "Invalid HTTP method. You can only use the methods from the following list: "                               // Invalid HTTP method message.
	InvalidContentTypeNotification = "Invalid content type. You can use only them from the following list: "                                     // Invalid content type message.
	InvalidHeaderFormat            = "Invalid header format"                                                                                     // Invalid header format message.
)

// Databases.
const (
	MongoDB = "MongoDB" // MongoDB database name.
)

// Domains
const (
	UseCase = "UseCaseV1" // UseCase domain name.
)

// Deliveries
const (
	Gin = "Gin" // Gin delivery name.
)

// Common routes
const (
	GetAllItemsURL = ""     // Endpoint for fetching all items.
	GetItemByIdURL = "/:id" // Endpoint for fetching an item by its ID.
	ItemIdParam    = "id"   // Parameter name for item ID.
)

// Domain routes.
const (
	UsersGroupPath = "/users" // Users domain route.
	// Initialize other routes here
)

// Database table names.
const (
	UsersTable = "users" // UsersTable represents the name of the users table in the database.
)

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string
