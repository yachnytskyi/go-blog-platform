package domain

// CommonValidator encapsulates validation rules for a field.
type CommonValidator struct {
	FieldName    string // The name of the field being validated.
	FieldRegex   string // The regular expression used to validate the field's characters.
	MinLength    int    // The minimum allowed length for the field.
	MaxLength    int    // The maximum allowed length for the field.
	Notification string // The notification message for validation errors.
}

// NewCommonValidator creates a new instance of CommonValidator.
// Parameters:
//
//	fieldName: The name of the field being validated.
//	fieldRegex: The regular expression used to validate the field's characters.
//	minLength: The minimum allowed length for the field.
//	maxLength: The maximum allowed length for the field.
//	notification: The notification message for validation errors.
//
// Returns:
//
//	A new instance of CommonValidator.
func NewCommonValidator(fieldName, fieldRegex string, minLength, maxLength int, notification string) *CommonValidator {
	return &CommonValidator{
		FieldName:    fieldName,
		FieldRegex:   fieldRegex,
		MinLength:    minLength,
		MaxLength:    maxLength,
		Notification: notification,
	}
}
