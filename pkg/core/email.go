package core

import (
	"fmt"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Str("to", to).Msg("Failed to send email")
	}
}
