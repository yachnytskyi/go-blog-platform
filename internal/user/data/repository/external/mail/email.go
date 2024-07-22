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
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"gopkg.in/gomail.v2"
)

const (
	parsingMessage = "parsing template..."
)

// parseTemplateDirectory walks through the specified directory and parses all template files.
func parseTemplateDirectory(logger applicationModel.Logger, location, templatePath string) commonModel.Result[*template.Template] {
	var paths []string

	// Walk through the directory and gather all file paths.
	filePathWalkError := filepath.Walk(templatePath, func(path string, info os.FileInfo, walkError error) error {
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

	logger.Info(domainError.NewInfoMessage(location+".parseTemplateDirectory", parsingMessage))
	if validator.IsError(filePathWalkError) {
		internalError := domainError.NewInternalError(location+".parseTemplateDirectory."+parsingMessage, filePathWalkError.Error())
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[*template.Template](internalError)
	}

	// Parse all collected template files.
	parseFiles, parseFilesError := template.ParseFiles(paths...)
	if validator.IsError(parseFilesError) {
		internalError := domainError.NewInternalError(location+".ParseFiles."+parsingMessage, parseFilesError.Error())
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[*template.Template](internalError)
	}

	return commonModel.NewResultOnSuccess[*template.Template](parseFiles)
}

// SendEmail sends an email to the specified user using the provided email data.
func SendEmail(logger applicationModel.Logger, location string, user userModel.User, data userModel.EmailData) error {
	emailConfig := config.GetEmailConfig()
	smtpPass := emailConfig.SMTPPassword
	smtpUser := emailConfig.SMTPUser
	smtpHost := emailConfig.SMTPHost
	smtpPort := emailConfig.SMTPPort

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	prepareSendMessage := prepareSendMessage(logger, location+".SendEmail", user.Email, data)
	if validator.IsError(prepareSendMessage.Error) {
		return prepareSendMessage.Error
	}

	dialAndSendError := dialer.DialAndSend(prepareSendMessage.Data)
	if validator.IsError(dialAndSendError) {
		internalError := domainError.NewInternalError(location+".SendEmail.DialAndSend", dialAndSendError.Error())
		logger.Error(internalError)
		return internalError
	}

	return nil
}

// prepareSendMessage prepares the email message to be sent.
func prepareSendMessage(logger applicationModel.Logger, location, userEmail string, data userModel.EmailData) commonModel.Result[*gomail.Message] {
	emailConfig := config.GetEmailConfig()
	from := emailConfig.EmailFrom
	to := userEmail

	// Parse the template directory to get the templates.
	var body bytes.Buffer
	template := parseTemplateDirectory(logger, location+".prepareSendMessage", data.TemplatePath)
	if validator.IsError(template.Error) {
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.parseTemplateDirectory", template.Error.Error())
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Retrieve the specific email template.
	emailTemplate := template.Data.Lookup(data.TemplateName)
	if emailTemplate == nil {
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.TemplateNotFound", constants.EmailTemplateNotFound)
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Execute the template to generate the email body.
	executeError := emailTemplate.Execute(&body, &data)
	if validator.IsError(executeError) {
		internalError := domainError.NewInternalError(location+".SendEmail.PrepareSendMessage.Execute", executeError.Error())
		logger.Error(internalError)
		return commonModel.NewResultOnFailure[*gomail.Message](internalError)
	}

	// Create a new email message.
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	return commonModel.NewResultOnSuccess[*gomail.Message](message)
}
