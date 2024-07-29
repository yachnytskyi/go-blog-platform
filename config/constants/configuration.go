package constants

import (
	"net/http"
	"time"
)

// Pagination constants.
const (
	Page                 = "page"       // Page query parameter.
	Limit                = "limit"      // Limit query parameter.
	OrderBy              = "order_by"   // Order by query parameter.
	SortOrder            = "sort_order" // Sort order query parameter.
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

// DateTime formats.
const (
	DateTimeFormat       = "02-Jan-2006 03:04:05 PM MST" // Format for human-readable dates.
	LoggerDateTimeFormat = time.RFC3339                  // Format for machine-readable dates (ISO 8601).
)

// Token-related constants.
const (
	AccessTokenValue                            = "access_token"                          // Access token value.
	RefreshTokenValue                           = "refresh_token"                         // Refresh token value.
	LoggedInValue                               = "logged_in"                             // Logged in status value.
	TokenDomainValue                            = "localhost"                             // Token domain value.
	LogoutMaxAgeValue                           = -1                                      // Logout max age value.
	User                             contextKey = "user"                                  // User context key.
	ID                               contextKey = "id"                                    // ID context key.
	UserRole                         contextKey = "userRole"                              // User role context key.
	IDContextMissing                            = "ID context value is missing or empty." // ID context missing error message.
	PasswordResetTokenExpirationTime            = time.Hour * 24                          // PasswordResetTokenExpirationTime represents the duration after which a password reset token expires.
)

// Email subjects and URLs.
const (
	EmailConfirmationUrl     = "users/verifyemail/"                                     // Email confirmation URL.
	ForgottenPasswordUrl     = "users/reset-password/"                                  // Forgotten password URL.
	EmailConfirmationSubject = "Your account verification code"                         // Email confirmation subject.
	ForgottenPasswordSubject = "Your password reset token (it is valid for 15 minutes)" // Forgotten password subject.
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
