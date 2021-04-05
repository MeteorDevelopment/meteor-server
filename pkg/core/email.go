package core

import (
	"fmt"
	"net/smtp"
)

var emailAuth smtp.Auth

func InitEmail() {
	emailAuth = smtp.PlainAuth("", "noreply@meteorclient.com", GetPrivateConfig().EmailPassword, "smtp.zoho.eu")
}

func SendEmail(to string, subject string, text string) {
	msg := []byte(fmt.Sprintf("From: Meteor <noreply@meteorclient.com>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", to, subject, text))
	err := smtp.SendMail("smtp.zoho.eu:587", emailAuth, "noreply@meteorclient.com", []string{to}, msg)

	if err != nil {
		fmt.Printf("[Email] Failed to send email to '%s'\n", to)
		fmt.Printf("        %s\n", err)
	}
}
