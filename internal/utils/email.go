// Package utils
package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Define dialer interface
type dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// Function to create dialer (can be mocked in tests)
var newDialer = func(host string, port int, username, password string) dialer {
	return gomail.NewDialer(host, port, username, password)
}

func SendVerificationEmail(recipient, verificationURL string) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	port := 587
	if p, err := strconv.Atoi(portStr); err == nil {
		port = p
	}

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Verify Your Email Address")
	m.SetBody("text/plain", fmt.Sprintf("Click the link to verify your email: %s", verificationURL))
	d := newDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send verification email to %s: %v", recipient, err)
	} else {
		log.Printf("Verification email sent to %s", recipient)
	}
}
