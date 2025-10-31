// Package services
package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
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
	return s.sendEmailWithRetry(to, subject, EmailTypeWelcome, data, 2)
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
	return s.sendEmailWithRetry(to, subject, EmailTypePasswordReset, data, 2)
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
	return s.sendEmailWithRetry(to, subject, EmailTypeLoginNotification, data, 2)
}

func (s *emailService) sendEmailWithRetry(to, subject string, emailType EmailType, data EmailData, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := s.sendEmail(to, subject, emailType, data); err != nil {
			lastErr = err
			log.Printf("Email send attempt %d failed for %s: %v", i+1, to, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

func (s *emailService) sendEmail(to, subject string, emailType EmailType, data EmailData) error {
	body, err := s.renderTemplate(emailType, data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

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

	addr := s.config.SMTPHost + ":" + strconv.Itoa(s.config.SMTPPort)

	// Try different connection methods based on port
	switch s.config.SMTPPort {
	case 587:
		// STARTTLS (most common for email providers)
		return s.sendWithStartTLS(addr, to, message)
	case 465:
		// TLS/SSL (direct TLS connection)
		return s.sendWithTLS(addr, to, message)
	case 25:
		// Plain SMTP (usually without TLS)
		return s.sendPlain(addr, to, message)
	default:
		// Auto-detect based on common ports
		if s.config.SMTPPort == 465 {
			return s.sendWithTLS(addr, to, message)
		}
		return s.sendWithStartTLS(addr, to, message)
	}
}

func (s *emailService) sendWithStartTLS(addr, to, message string) error {
	// Connect to SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("SMTP dial failed: %w", err)
	}
	defer client.Close()

	// Send STARTTLS command
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         s.config.SMTPHost,
			InsecureSkipVerify: false,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	// Send email
	if err = client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("SMTP MAIL failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	defer w.Close()

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	return nil
}

func (s *emailService) sendWithTLS(addr, to, message string) error {
	tlsConfig := &tls.Config{
		ServerName:         s.config.SMTPHost,
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS connection failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("SMTP MAIL failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	defer w.Close()

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	return nil
}

func (s *emailService) sendPlain(addr, to, message string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(message))
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
