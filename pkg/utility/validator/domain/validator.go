package domain

// CommonValidator encapsulates validation rules for a field.
type CommonValidator struct {
	FieldName    string // The name of the field being validated.
	FieldRegex   string // The regular expression used to validate the field's characters.
	MinLength    int    // The minimum allowed length for the field.
	MaxLength    int    // The maximum allowed length for the field.
	Notification string // The notification message for validation errors.
}
