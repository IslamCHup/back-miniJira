package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type EmailService struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewEmailService() *EmailService {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
	
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

	msg := []byte(
		"From: " + s.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
			body + "\r\n",
	)

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial failed: %w", err)
	}
	defer c.Quit()

	if ok, _ := c.Extension("STARTTLS"); !ok {
		return fmt.Errorf("server does not support STARTTLS")
	}
	if err := c.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
		return fmt.Errorf("starttls failed: %w", err)
	}

	if ok, _ := c.Extension("AUTH"); !ok {
		return fmt.Errorf("server does not support AUTH")
	}
	if err := c.Auth(smtp.PlainAuth("", s.user, s.pass, s.host)); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	if err := c.Mail(s.user); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA open failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}

	return nil
}

func (s *EmailService) SendVerificationEmail(to, name, link string) error {
	subject := "Подтверждение аккаунта MiniJira"
	body := fmt.Sprintf("Здравствуйте, %s!\n\nПерейдите по ссылке для подтверждения:\n%s", name, link)
	return s.SendEmail(to, subject, body)
}
