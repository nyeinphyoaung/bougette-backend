package utilities

import (
	"bougette-backend/configs"
	"bytes"
	"embed"
	"path/filepath"
	"strconv"
	"text/template"

	"gopkg.in/gomail.v2"
)

//go:embed  templates
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string
}

type MailData struct {
	Name    string
	Subject string
	Meta    interface{}
}

func NewMailer() Mailer {
	mailSender := configs.Envs.MAIL_SENDER
	mailHost := configs.Envs.MAIL_HOST
	mailPort, _ := strconv.Atoi(configs.Envs.MAIL_PORT)
	mailUsername := configs.Envs.MAIL_USERNAME
	mailPassword := configs.Envs.MAIL_PASSWORD

	dialer := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)
	return Mailer{
		dialer: dialer,
		sender: mailSender,
	}
}

func (mailer *Mailer) SendViaMail(recipient string, templateFile string, data MailData) error {
	absolutePath := filepath.Join("templates", templateFile)
	tmpl, err := template.ParseFS(templateFS, absolutePath)

	if err != nil {
		return err
	}

	data.Name = configs.Envs.VIA_APP_NAME

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	gomailMsg := gomail.NewMessage()
	gomailMsg.SetHeader("To", recipient)
	gomailMsg.SetHeader("From", mailer.sender)
	gomailMsg.SetHeader("Subject", subject.String())
	gomailMsg.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(gomailMsg)
	if err != nil {
		return err
	}

	return nil
}
