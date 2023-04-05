package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// Email template parser.
func SendEmail(user *models.UserFullResponse, data *EmailData, temp *template.Template, templateName string) error {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	// Send data.
	from := config.EmailFrom
	smtpPass := config.SMTPPassword
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	if err := temp.ExecuteTemplate(&body, templateName, &data); err != nil {
		log.Fatal("Could not execute template", err)
	}

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
		return err
	}
	return nil
}
