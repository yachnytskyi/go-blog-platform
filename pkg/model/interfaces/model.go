package interfaces

// ServerRouters holds the routers for different modules of the application.
type ServerRouters struct {
	UserRouter Router
	PostRouter Router
	// Add other routers as needed.
}

// NewServerRouters creates a new instance of ServerRouters with the given routers.
func NewServerRouters(userRouter Router, postRouter Router) ServerRouters {
	return ServerRouters{
		UserRouter: userRouter,
		PostRouter: postRouter,
		// Add other routers as needed.
	}
}

// EmailData holds the data required for sending an email.
type EmailData struct {
	Recipient    string // Recipient's email address.
	URL          string // URL to be included in the email.
	TemplateName string // Name of the email template.
	TemplatePath string // Path to the email template.
	FirstName    string // Recipient's first name.
	Subject      string // Subject of the email.
}

func NewEmailData(recipient, url, templateName, templatePath, firstName, subject string) EmailData {
	return EmailData{
		Recipient:    recipient,
		URL:          url,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    firstName,
		Subject:      subject,
	}
}
