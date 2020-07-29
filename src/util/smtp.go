package util

import (
	"tree-hole/config"

	"github.com/go-gomail/gomail"
)

var (
	smtpDialer   *gomail.Dialer
	emailAddress string
)

func initUtilSMTP() {
	smtp := config.Config.SMTP
	smtpDialer = gomail.NewDialer(smtp.SMTPAddress, smtp.SMTPPort, smtp.EmailAddress, smtp.EmailPassword)
}

func SendEmail(receiver string, subject string, content string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", emailAddress)
	mail.SetHeader("To", receiver)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", content)
	return smtpDialer.DialAndSend(mail)
}
