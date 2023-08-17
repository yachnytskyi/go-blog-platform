package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	"gopkg.in/gomail.v2"
)

// Email template parser.
func ParseTemplateDirectory(directory string) (*template.Template, error) {
	var paths []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			var sendEmailInternalError *domain_error.InternalError = new(domain_error.InternalError)
			sendEmailInternalError.Location = "User.Data.Repository.External.Mail.ParseTemplateDirectory.Walk"
			sendEmailInternalError.Reason = err.Error()
			fmt.Println(sendEmailInternalError)
			return sendEmailInternalError
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("parsing template...")

	if err != nil {
		var sendEmailInternalError *domain_error.InternalError = new(domain_error.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.parsing template"
		sendEmailInternalError.Reason = err.Error()
		fmt.Println(sendEmailInternalError)
		return nil, sendEmailInternalError
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *userModel.User, data *userModel.EmailData, templateName string) error {
	config, err := config.LoadConfig(".")

	if err != nil {
		var sendEmailInternalError *domain_error.InternalError = new(domain_error.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.LoadConfig"
		sendEmailInternalError.Reason = err.Error()
		fmt.Println(sendEmailInternalError)
		return sendEmailInternalError
	}

	// Send data.
	from := config.EmailFrom
	smtpPass := config.SMTPPassword
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer
	template, err := ParseTemplateDirectory("internal/user/delivery/http/utility/template")

	if err != nil {
		var sendEmailInternalError *domain_error.InternalError = new(domain_error.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.ParseTemplateDirectory"
		sendEmailInternalError.Reason = "could not parse template"
		fmt.Println(sendEmailInternalError)
		return sendEmailInternalError
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	if err := dialer.DialAndSend(message); err != nil {
		var sendEmailInternalError *domain_error.InternalError = new(domain_error.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.DialAndSend"
		sendEmailInternalError.Reason = err.Error()
		fmt.Println(sendEmailInternalError)
		return sendEmailInternalError
	}

	return nil
}

func UserFirstName(firstName string) string {
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	return firstName
}
