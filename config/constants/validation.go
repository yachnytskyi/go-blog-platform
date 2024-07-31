package constants

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
