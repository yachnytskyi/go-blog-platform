package constants

// contextKey is a custom type for context keys to prevent collisions.
type contextKey string

// Application version.
const (
	AppVersion = V1   // AppVersion represents the application's current version.
	V1         = "v1" // V1 specifies the version value of the application.
)

// Environment configuration paths.
const (
	DefaultEnvironmentsPath = EnvironmentsPath           // Default path for environment configuration files.
	EnvironmentsPath        = "config/environment/.env." // Base path for environment configuration files.

	// Environment types.
	DefaultEnvironment       = LocalEnvironment // Default environment defines the default environment setting for the application.
	Environment              = LocalEnvironment // Environment defines the environment setting for the application.
	TestEnvironment          = "test"           // Test development environment.
	LocalEnvironment         = "local"          // Local development environment.
	DockerDevEnvironment     = "docker.dev"     // Docker development environment.
	DockerStagingEnvironment = "docker.staging" // Docker staging environment.
	DockerProductEnvironment = "docker.prod"    // Docker productuction environment.

	// Configuration paths.
	DefaultConfigPath = YamlConfigPath                          // Default path for the application configuration file.
	ConfigPath        = YamlConfigPath                          // Path for the application configuration file.
	YamlConfigPath    = "config/yaml/v1/local.application.yaml" // Path for the application yaml configuration file.

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
