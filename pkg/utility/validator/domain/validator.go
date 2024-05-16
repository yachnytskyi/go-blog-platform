package domain

type CommonValidator struct {
	FieldName    string
	FieldRegex   string
	MinLength    int
	MaxLength    int
	Notification string
}
