package service

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) RegisterNotify(email string) error {
	sender := os.Getenv("EMAILSENDER")
	smtp := os.Getenv("SMTP")
	smtpport := os.Getenv("SMTPPORT")
	port, err := strconv.Atoi(smtpport)
	if err != nil {
		return fmt.Errorf("fail to convert smtpport to int")
	}
	emailPassword := os.Getenv("EMAILPASSWORD")
	subject := "Registration Successful!"
	text := "Welcome to our website,<br><p>we hope you have great experience on our shopping sotre<p>Best regards"

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(smtp, port, sender, emailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("fail to dial and send:%w", err)
	}
	return nil
}
func (s *EmailService) SetResetPasswordLink(email, encryptEmail string) error {
	domainname := os.Getenv("RESETLINK")
	if domainname == "" {
		fmt.Println("RESETLINK environment variable not set")
	}
	link := fmt.Sprintf("%s/%s", domainname, encryptEmail)

	sender := os.Getenv("EMAILSENDER")
	smtp := os.Getenv("SMTP")
	smtpport := os.Getenv("SMTPPORT")
	port, err := strconv.Atoi(smtpport)
	if err != nil {
		return fmt.Errorf("fail to convert smtpport to int")
	}
	emailPassword := os.Getenv("EMAILPASSWORD")
	subject := "Link to Reset Password"

	text := fmt.Sprintf("Hello, Click the link here to reset your password, You gonna received a new Password after clicking this<br><a>%s<a>", link)

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(smtp, port, sender, emailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("fail to dial and send:%w", err)
	}
	return nil
}
func (s *EmailService) NewPasswordNotify(decryptEmail, newPassword string) error {

	sender := os.Getenv("EMAILSENDER")
	smtp := os.Getenv("SMTP")
	smtpport := os.Getenv("SMTPPORT")
	port, err := strconv.Atoi(smtpport)
	if err != nil {
		return fmt.Errorf("fail to convert smtpport to int")
	}
	emailPassword := os.Getenv("EMAILPASSWORD")
	subject := "Reset Password"

	text := fmt.Sprintf("Hello, System generate a temporary password for you, Sign in to change your password as you wish<br><p>%s<p>", newPassword)

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", decryptEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(smtp, port, sender, emailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("fail to dial and send:%w", err)
	}
	return nil
}
