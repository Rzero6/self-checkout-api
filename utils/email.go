package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string, attachments []string) error {
	appName := os.Getenv("APP_NAME")
	from := os.Getenv("APP_EMAIL_ADDRESS")
	password := os.Getenv("APP_EMAIL_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", appName+"<noreply@"+appName+".com>")
	m.SetHeader("Reply-to", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	for _, attachment := range attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	return d.DialAndSend(m)
}
