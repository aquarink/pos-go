package utils

import (
	"net/smtp"
	"os"
)

func SendEmail(to, subject, text, html string) error {
	from := os.Getenv("EMAIL_FROM")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Setup email message
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + html

	// Authentication
	auth := smtp.PlainAuth("", username, password, host)

	// Sending email
	err := smtp.SendMail(host+":"+port, auth, username, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
