package constants

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string

// Application version.
const (
	Version = "v1" // Version of the application.
)

// Environment configuration paths.
const (
	EnvironmentsPath                   = "config/environment/.env."                  // Base path for environment configuration files.
	LocalDevEnvironment                = "local.dev"                                 // Local development environment.
	DockerDevEnvironment               = "docker.dev"                                // Docker development environment.
	DefaultEnvironmentsPath            = "config/environment/.env."                  // Default path for environment configuration files.
	DefaultConfigPath                  = "config/yaml/v1/local.dev.application.yaml" // Default path for the application configuration file.
	DefaultConfigPathNotification      = "Using default configuration path"          // 	// Notification message for using the default configuration path.
	DefaultEnrironmentPathNotification = "Using default environment path"            // 	// Notification message for using the default environment path.
)

// Configuration libraries used in the application.
const (
	Config = Viper   // Configuration library to use.
	Viper  = "Viper" // Viper configuration library.
)

// Logger libraries used in the application.
const (
	Zerolog = "Zerolog" // Zerolog logger library.
)

// Databases used in the application.
const (
	MongoDB = "MongoDB" // MongoDB database name.
)

// Domains used in the application.
const (
	UseCaseV1 = "UseCaseV1" // UseCaseV1 domain name.
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
	PostsGroupPath = "/posts" // Posts domain route.
	// Initialize other routes here.
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
)

// Schemes used in the application.
const (
	HTTP  = "http"  // HTTP scheme.
	HTTPS = "https" // HTTPS scheme.
)

// HTTP Headers used in the application.
const (
	CorrelationIDHeader = "X-Correlation-ID" // Correlation-ID header for tracking requests across systems.
	ContentType         = "Content-Type"     // Content-Type header.
	Authorization       = "Authorization"    // Authorization header.
	Bearer              = "Bearer"           // Bearer token prefix.
)
