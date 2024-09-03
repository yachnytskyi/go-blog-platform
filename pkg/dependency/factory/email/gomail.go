package email

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"gopkg.in/gomail.v2"
)

const (
	parsingMessage     = "parsing template..."
	messageFrom        = "From"
	messageTo          = "To"
	messageHeader      = "Subject"
	messageBody        = "text/html"
	messageAlternative = "text/plain"
)

type GoMail struct {
	Config *config.ApplicationConfig
	Logger interfaces.Logger
}

func NewGoMail(config *config.ApplicationConfig, logger interfaces.Logger) GoMail {
	return GoMail{
		Config: config,
		Logger: logger,
	}
}

// SendEmail sends an email to the specified user using the provided email data.
func (goMail GoMail) SendEmail(config *config.ApplicationConfig, logger interfaces.Logger, location string, data any, emailData interfaces.EmailData) error {
	smtpPass := config.Email.SMTPPassword
	smtpUser := config.Email.SMTPUser
	smtpHost := config.Email.SMTPHost
	smtpPort := config.Email.SMTPPort

	prepareSendMessage := prepareSendMessage(config, logger, location+".SendEmail", emailData)
	if validator.IsError(prepareSendMessage.Error) {
		return prepareSendMessage.Error
	}

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	dialAndSendError := dialer.DialAndSend(prepareSendMessage.Data)
	if validator.IsError(dialAndSendError) {
		internalError := domainError.NewInternalError(location+".SendEmail.DialAndSend", dialAndSendError.Error())
		logger.Error(internalError)
		return internalError
	}

	return nil
}

// prepareSendMessage prepares the email message to be sent.
func prepareSendMessage(config *config.ApplicationConfig, logger interfaces.Logger, location string, data interfaces.EmailData) common.Result[*gomail.Message] {
	from := config.Email.EmailFrom

	// Parse the template directory to get the templates.
	var body bytes.Buffer
	template := parseTemplateDirectory(logger, location+".prepareSendMessage", data.TemplatePath)
	if validator.IsError(template.Error) {
		internalError := domainError.NewInternalError(location+".prepareSendMessage.parseTemplateDirectory", template.Error.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Retrieve the specific email template.
	emailTemplate := template.Data.Lookup(data.TemplateName)
	if emailTemplate == nil {
		internalError := domainError.NewInternalError(location+".prepareSendMessage.Lookup", constants.EmailTemplateNotFound)
		logger.Error(internalError)
		return common.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Execute the template to generate the email body.
	executeError := emailTemplate.Execute(&body, &data)
	if validator.IsError(executeError) {
		internalError := domainError.NewInternalError(location+".prepareSendMessage.Execute", executeError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Create a new email message.
	message := gomail.NewMessage()
	message.SetHeader(messageFrom, from)
	message.SetHeader(messageTo, data.Recipient)
	message.SetHeader(messageHeader, data.Subject)
	message.SetBody(messageBody, body.String())
	message.AddAlternative(messageAlternative, html2text.HTML2Text(body.String()))
	return common.NewResultOnSuccess[*gomail.Message](message)
}

// parseTemplateDirectory walks through the specified directory and parses all template files.
func parseTemplateDirectory(logger interfaces.Logger, location, templatePath string) common.Result[*template.Template] {
	var paths []string

	// Walk through the directory and gather all file paths.
	walkError := filepath.Walk(templatePath, func(path string, info os.FileInfo, walkError error) error {
		if validator.IsError(walkError) {
			internalError := domainError.NewInternalError(location+".parseTemplateDirectory.Walk", walkError.Error())
			logger.Error(internalError)
			return internalError
		}
		if info.IsDir() {
			return nil // Skip directories.
		}

		paths = append(paths, path) // Collect file paths.
		return nil
	})

	logger.Debug(domainError.NewInfoMessage(location+".parseTemplateDirectory", parsingMessage))
	if validator.IsError(walkError) {
		internalError := domainError.NewInternalError(location+".parseTemplateDirectory."+parsingMessage, walkError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*template.Template](internalError)
	}

	// Parse all collected template files.
	parseFiles, parseFilesError := template.ParseFiles(paths...)
	if validator.IsError(parseFilesError) {
		internalError := domainError.NewInternalError(location+".ParseFiles."+parsingMessage, parseFilesError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*template.Template](internalError)
	}

	return common.NewResultOnSuccess[*template.Template](parseFiles)
}
