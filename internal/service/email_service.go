package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("SMTP_FROM"),
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	msg := []byte(
		"From: " + s.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body,
	)

	return smtp.SendMail(addr, auth, s.user, []string{to}, msg)
}

func (s *EmailService) SendVerificationEmail(to, name, link string) error {
	subject := "Подтверждение аккаунта MiniJira"
	body := fmt.Sprintf("Здравствуйте, %s!\n\nПерейдите по ссылке для подтверждения:\n%s", name, link)
	return s.SendEmail(to, subject, body)
}
