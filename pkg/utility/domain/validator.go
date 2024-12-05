package domain

import "regexp"

// StringValidator defines validation rules for string fields.
type StringValidator struct {
	FieldName    string         // The name of the field being validated.
	Field        string         // The value of the field being validated.
	FieldRegex   *regexp.Regexp // The compiled regular expression used to validate the field's characters.
	MinLength    int            // The minimum allowed length for the field.
	MaxLength    int            // The maximum allowed length for the field.
	IsOptional   bool           // Indicates if the field is optional.
}

func NewStringValidator(fieldName string, field string, fieldRegex *regexp.Regexp, minLength, maxLength int,  isOptional bool) StringValidator {
	return StringValidator{
		FieldName:    fieldName,
		Field:        field,
		FieldRegex:   fieldRegex,
		MinLength:    minLength,
		MaxLength:    maxLength,
		IsOptional:   isOptional,
	}
}
