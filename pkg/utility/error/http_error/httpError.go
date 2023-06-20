package http_error

import "fmt"

type HttpValidationError struct {
	Field     string `json:"field"`
	FieldType string `json:"type"`
	Reason    string `json:"reason"`
}

func NewHttpValidationError(field string, fieldType string, reason string) error {
	return HttpValidationError{
		Field:     field,
		FieldType: fieldType,
		Reason:    reason,
	}
}

func (err HttpValidationError) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " reason: " + err.Reason)
}
