package constants

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string

// Application version.
const (
	version = "version" // Version specifies the constant value of the application.
)

// Environment configuration paths.
const (
	// Configuration paths.
	EnvironmentsPath        = "config/environment/.env"               // Base path for environment configuration files.
	DefaultEnvironmentsPath = "config/environment/.env"               // Default path for environment configuration files.
	DefaultEnvironment      = ".local"                                // Default environment defines the default environment setting for the application.
	DefaultConfigPath       = "config/yaml/v1/local.application.yaml" // Default path for the application configuration file.

	// Notification messages.
	DefaultConfigPathNotification      = "Using default configuration path" // Notification message for using the default configuration path.
	DefaultEnvironmentPathNotification = "Using default environment path"   // Notification message for using the default environment path.
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

// Email libraries used in the application.
const (
	GoMail = "GoMail" // GoMail email library.
)

// Databases used in the application.
const (
	MongoDB = "MongoDB" // MongoDB database name.
)

// Deliveries used in the application.
const (
	Gin = "Gin" // Gin delivery name.
)
