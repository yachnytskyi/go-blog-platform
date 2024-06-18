// Package mail provides functionality for sending emails using SMTP.
package mail

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"gopkg.in/gomail.v2"
)

const (
	loggerMessage = "parsing template..."
)

// ParseTemplateDirectory walks through the specified directory and parses all template files.
// It returns a commonModel.Result containing a pointer to the parsed template.
// Parameters:
// - location: a string used to identify the source location in error messages and logging.
// - templatePath: the path to the directory containing the template files.
func ParseTemplateDirectory(location, templatePath string) commonModel.Result[*template.Template] {
	var paths []string

	// Walk through the directory and gather all file paths.
	filePathWalkError := filepath.Walk(templatePath, func(path string, info os.FileInfo, walkError error) error {
		if validator.IsError(walkError) {
			// Handle any error encountered during walking the directory.
			internalError := domainError.NewInternalError(location+".ParseTemplateDirectory.Walk", walkError.Error())
			logging.Logger(internalError)
			return internalError
		}
		if info.IsDir() {
			return nil // Skip directories.
		}

		paths = append(paths, path) // Collect file paths.
		return nil
	})

	logging.Logger(loggerMessage)
	if validator.IsError(filePathWalkError) {
		// Handle error if walking the directory fails.
		internalError := domainError.NewInternalError(location+".ParseTemplateDirectory."+loggerMessage, filePathWalkError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*template.Template](internalError)
	}

	// Parse all collected template files.
	parseFiles, parseFilesError := template.ParseFiles(paths...)
	if validator.IsError(parseFilesError) {
		// Handle error if parsing the templates fails.
		internalError := domainError.NewInternalError(location+".ParseFiles."+loggerMessage, parseFilesError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*template.Template](internalError)
	}

	// Return the successfully parsed templates.
	return commonModel.NewResultOnSuccess[*template.Template](parseFiles)
}

// SendEmail sends an email to the specified user using the provided email data.
// It returns an error if sending the email fails.
// Parameters:
// - location: a string used to identify the source location in error messages and logging.
// - user: the recipient of the email.
// - data: the data to be used in the email body.
func SendEmail(location string, user userModel.User, data userModel.EmailData) error {
	// Email configuration.
	emailConfig := config.GetEmailConfig()
	smtpPass := emailConfig.SMTPPassword
	smtpUser := emailConfig.SMTPUser
	smtpHost := emailConfig.SMTPHost
	smtpPort := emailConfig.SMTPPort

	// Create a new dialer with the email configuration.
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Prepare the email message.
	prepareSendMessage := prepareSendMessage(location+".SendEmail", user.Email, data)
	if validator.IsError(prepareSendMessage.Error) {
		// Handle error if preparing the email message fails.
		return prepareSendMessage.Error
	}

	// Send the email.
	dialAndSendError := dialer.DialAndSend(prepareSendMessage.Data)
	if validator.IsError(dialAndSendError) {
		// Handle error if sending the email fails.
		internalError := domainError.NewInternalError(location+".SendEmail.DialAndSend", dialAndSendError.Error())
		logging.Logger(internalError)
		return internalError
	}

	return nil // Return nil if email is sent successfully.
}

// prepareSendMessage prepares the email message to be sent.
// It returns a commonModel.Result containing a pointer to the prepared message.
// Parameters:
// - location: a string used to identify the source location in error messages and logging.
// - userEmail: the recipient's email address.
// - data: the data to be used in the email body, including the template path and name.
func prepareSendMessage(location, userEmail string, data userModel.EmailData) commonModel.Result[*gomail.Message] {
	// Load email configuration.
	emailConfig := config.GetEmailConfig()

	// Set sender and recipient.
	from := emailConfig.EmailFrom
	to := userEmail

	// Prepare the email body buffer.
	var body bytes.Buffer

	// Parse the template directory to get the templates.
	template := ParseTemplateDirectory(location+".prepareSendMessage", data.TemplatePath)
	if validator.IsError(template.Error) {
		// Handle template parsing error.
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.ParseTemplateDirectory", template.Error.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Retrieve the specific email template.
	emailTemplate := template.Data.Lookup(data.TemplateName)
	if emailTemplate == nil {
		// Handle missing email template.
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.TemplateNotFound", constants.EmailTemplateNotFound)
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Execute the template to generate the email body.
	executeError := emailTemplate.Execute(&body, &data)
	if validator.IsError(executeError) {
		// Handle template execution error.
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.Execute", executeError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Create a new email message.
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	// Return the prepared email message.
	return commonModel.NewResultOnSuccess[*gomail.Message](message)
}
