package constants

// Regex patterns and string length constraints.
const (
	DefaultStringRegex             = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`     // Default string regex pattern.
	DefaultTitleStringRegex        = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`                  // Default title string regex pattern.
	DefaultTextStringRegex         = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$` // Default text string regex pattern.
	DefaultMinStringLength         = 4                                                      // Default minimum string length.
	DefaultMaxStringLength         = 40                                                     // Default maximum string length.
	DefaultMinOptionalStringLength = 0                                                      // Default minimum optional string length.
	DefaultMaxOptionalStringLength = 40                                                     // Default maximum optional string length.
	FieldRequired                  = "required"                                             // Field required status.
	FieldOptional                  = "optional"                                             // Field optional status.
	True                           = "true"                                                 // True value.
	False                          = "false"                                                // False value.
)
