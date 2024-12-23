package interfaces

type ServerRouters struct {
	HealthCheckRouter Router
	UserRouter        Router
	PostRouter        Router
	// Add other routers as needed.
}

func NewServerRouters(healthCheckRouter, userRouter, postRouter Router) ServerRouters {
	return ServerRouters{
		HealthCheckRouter: healthCheckRouter,
		UserRouter:        userRouter,
		PostRouter:        postRouter,
		// Add other routers as needed.
	}
}

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
