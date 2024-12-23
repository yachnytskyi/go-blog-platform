package constants

import (
	"time"
)

// Pagination constants.
const (
	Page                     = "page"       // Page query parameter.
	Limit                    = "limit"      // Limit query parameter.
	OrderBy                  = "order_by"   // Order by query parameter.
	SortOrder                = "sort_order" // Sort order query parameter.
	DefaultPage              = "1"          // Default page number.
	DefaultLimit             = "10"         // Default limit per page.
	DefaultAmountOfPageLinks = 10           // Default amount of page links for pagination.
	DefaultPageInteger       = 1            // Default page number.
	DefaultLimitInteger      = 10           // Default limit per page.
	MaxItemsPerPage          = 100          // Maximum items per page.
	DefaultSortOrder         = SortDescend  // Default sort order.
	DefaultOrderBy           = "created_at" // Default order by field.
	SortAscend               = "ascend"     // Ascending sort order.
	SortDescend              = "descend"    // Descending sort order.
)

// Context timers.
const (
	DefaultContextTimer = 5 * time.Second // Default context timeout duration.
)

// DateTime formats.
const (
	DateTimeFormat       = "02-Jan-2006 03:04:05 PM MST" // Format for human-readable dates.
	LoggerDateTimeFormat = time.RFC3339                  // Format for machine-readable dates (ISO 8601).
)

// Error message formatting constants.
const (
	BaseErrorMessageFormat = "location: %s notification: %s" // Format string for logging or displaying error messages with location and notification details.
)

// Token-related constants.
const (
	AccessTokenValue                            = "access_token"                          // Access token value.
	RefreshTokenValue                           = "refresh_token"                         // Refresh token value.
	LoggedInValue                               = "logged_in"                             // Logged in status value.
	LogoutMaxAgeValue                           = -1                                      // Logout max age value.
	User                             contextKey = "user"                                  // User context key.
	ID                               contextKey = "id"                                    // ID context key.
	UserRole                         contextKey = "userRole"                              // User role context key.
	IDContextMissing                            = "ID context value is missing or empty." // ID context missing error message.
	PasswordResetTokenExpirationTime            = time.Hour * 24                          // PasswordResetTokenExpirationTime represents the duration after which a password reset token expires.
)

// Common routes used in the application.
const (
	GetAllItemsURL = ""     // Endpoint for fetching all items.
	GetItemByIdURL = "/:id" // Endpoint for fetching an item by its ID.
	ItemIdParam    = "id"   // Parameter name for item ID.
)

// Domain-specific routes.
const (
	HealthGroupPath = "/health" // Health domain route.
	UsersGroupPath  = "/users"  // Users domain route.
	PostsGroupPath  = "/posts"  // Posts domain route.
	// Initialize other routes here.
)

// Email subjects and URLs.
const (
	EmailConfirmationUrl = "users/verifyemail/"    // Email confirmation URL.
	ForgottenPasswordUrl = "users/reset-password/" // Forgotten password URL.
)

// User route paths.
const (
	RegisterPath          = "/register"           // Registration route path.
	ForgottenPasswordPath = "/forgotten-password" // Forgotten password route path.
	ResetPasswordPath     = "/reset-password/:id" // Reset password route path with token.
	LoginPath             = "/login"              // Login route path.
	GetCurrentUserPath    = "/current_user"       // Get current user route path.
	UpdateCurrentUserPath = "/update"             // Update current user route path.
	DeleteCurrentUserPath = "/delete"             // Delete current user route path.
	RefreshTokenPath      = "/refresh"            // Refresh token route path.
	LogoutPath            = "/logout"             // Logout route path.
)

// Database table names.
const (
	UsersTable = "users" // Users table name in the database.
	PostsTable = "posts" // Posts table name in the database.
)

// Schemes used in the application.
const (
	HTTP  = "http"  // HTTP scheme.
	HTTPS = "https" // HTTPS scheme.
)

// HTTP Headers used in the application.
const (
	RequestIDHeader = "X-Request-ID"  // Request-ID header for tracking requests across systems.
	ContentType     = "Content-Type"  // Content-Type header.
	Authorization   = "Authorization" // Authorization header.
	Bearer          = "Bearer "       // Bearer token prefix.
)

// Operation status messages.
const (
	Success = "success" // Success message for successful operations.
	Fail    = "fail"    // Failure message for unsuccessful operations.
)

// HTTP methods for RESTful operations.
const (
	Get    = "GET"    // GET method for retrieving resources.
	Post   = "POST"   // POST method for creating resources.
	Put    = "PUT"    // PUT method for updating resources.
	Patch  = "PATCH"  // PATCH method for partially updating resources.
	Delete = "DELETE" // DELETE method for deleting resources.
)
