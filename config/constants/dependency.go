package constants

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string

// Application version.
const (
	Version = "v1" // Version of the application.
)

// Environment configuration paths.
const (
	EnvironmentsPath = "config/environment/.env." // Base path for environment configuration files.

	// Environment types
	TestEnvironment          = "test"        // Test development environment.
	LocalEnvironment         = "local"       // Local development environment.
	DockerDevEnvironment     = "docker.dev"  // Docker development environment.
	DockerProductEnvironment = "docker.prod" // Docker product environment.

	// Default configuration paths
	DefaultEnvironmentsPath = "config/environment/.env."                  // Default path for environment configuration files.
	DefaultConfigPath       = "config/yaml/v1/local.dev.application.yaml" // Default path for the application configuration file.

	// Notification messages
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
