package validator

// StringValidator defines validation rules for string fields.
type StringValidator struct {
	FieldName    string // The name of the field being validated.
	FieldRegex   string // The regular expression used to validate the field's characters.
	MinLength    int    // The minimum allowed length for the field.
	MaxLength    int    // The maximum allowed length for the field.
	Notification string // The notification message for validation errors.
	IsOptional   bool   // Indicates if the field is optional.
}

func NewStringValidator(fieldName, fieldRegex string, minLength, maxLength int, notification string) StringValidator {
	return StringValidator{
		FieldName:    fieldName,
		FieldRegex:   fieldRegex,
		MinLength:    minLength,
		MaxLength:    maxLength,
		Notification: notification,
	}
}
