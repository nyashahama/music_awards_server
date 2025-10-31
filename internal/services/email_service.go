// Package services
package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/templates"
)

type EmailType string

const (
	EmailTypeWelcome           EmailType = "welcome"
	EmailTypePasswordReset     EmailType = "password_reset"
	EmailTypeLoginNotification EmailType = "login_notification"
)

type EmailData struct {
	Username    string
	Token       string
	ResetURL    string
	AppURL      string
	LoginTime   string
	UserAgent   string
	IPAddress   string
	VoteCount   int
	CurrentYear int
}

type EmailService interface {
	SendWelcomeEmail(to, username string, voteCount int) error
	SendPasswordResetEmail(to, username, token string) error
	SendLoginNotificationEmail(to, username, userAgent, ipAddress string) error
}

type emailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) EmailService {
	return &emailService{config: cfg}
}

func (s *emailService) SendWelcomeEmail(to, username string, voteCount int) error {
	data := EmailData{
		Username:    username,
		VoteCount:   voteCount,
		AppURL:      s.getAppURL(),
		CurrentYear: time.Now().Year(),
	}

	subject := "Welcome to Music Awards!"
	body, err := s.renderTemplate(EmailTypeWelcome, data)
	if err != nil {
		return fmt.Errorf("failed to render welcome email template: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *emailService) SendPasswordResetEmail(to, username, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.getAppURL(), token)
	data := EmailData{
		Username:    username,
		Token:       token,
		ResetURL:    resetURL,
		CurrentYear: time.Now().Year(),
	}

	subject := "Reset Your Music Awards Password"
	body, err := s.renderTemplate(EmailTypePasswordReset, data)
	if err != nil {
		return fmt.Errorf("failed to render password reset email template: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *emailService) SendLoginNotificationEmail(to, username, userAgent, ipAddress string) error {
	data := EmailData{
		Username:    username,
		LoginTime:   time.Now().Format("January 2, 2006 at 3:04 PM MST"),
		UserAgent:   s.truncateUserAgent(userAgent),
		IPAddress:   ipAddress,
		CurrentYear: time.Now().Year(),
	}

	subject := "New Login to Your Music Awards Account"
	body, err := s.renderTemplate(EmailTypeLoginNotification, data)
	if err != nil {
		return fmt.Errorf("failed to render login notification email template: %w", err)
	}

	return s.sendEmail(to, subject, body)
}

func (s *emailService) renderTemplate(emailType EmailType, data EmailData) (string, error) {
	var tmplStr string

	switch emailType {
	case EmailTypeWelcome:
		tmplStr = templates.WelcomeEmailTemplate
	case EmailTypePasswordReset:
		tmplStr = templates.PasswordResetTemplate
	case EmailTypeLoginNotification:
		tmplStr = templates.LoginNotificationTemplate
	default:
		return "", fmt.Errorf("unknown email type: %s", emailType)
	}

	tmpl, err := template.New("email").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// SMTP authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// SMTP connection
	addr := s.config.SMTPHost + ":" + strconv.Itoa(s.config.SMTPPort)

	if s.config.UseTLS {
		// TLS connection
		tlsconfig := &tls.Config{
			ServerName: s.config.SMTPHost,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("TLS connection failed: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("SMTP client creation failed: %w", err)
		}
		defer client.Close()

		// Auth
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}

		// Set sender and recipient
		if err = client.Mail(s.config.FromEmail); err != nil {
			return fmt.Errorf("SMTP MAIL failed: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("SMTP RCPT failed: %w", err)
		}

		// Send email
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("SMTP DATA failed: %w", err)
		}
		defer w.Close()

		if _, err = w.Write([]byte(message)); err != nil {
			return fmt.Errorf("SMTP write failed: %w", err)
		}
	} else {
		// Plain SMTP (not recommended for production)
		err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(message))
		if err != nil {
			return fmt.Errorf("SMTP send failed: %w", err)
		}
	}

	return nil
}

func (s *emailService) getAppURL() string {
	// In production, this should be your frontend URL
	if url := os.Getenv("FRONTEND_URL"); url != "" {
		return url
	}
	return "https://music-awards.com"
}

func (s *emailService) truncateUserAgent(userAgent string) string {
	if len(userAgent) > 100 {
		return userAgent[:100] + "..."
	}
	return userAgent
}
