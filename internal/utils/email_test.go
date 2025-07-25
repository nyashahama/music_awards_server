package utils

import (
	"os"
	"testing"

	"gopkg.in/gomail.v2"
)

// Define dialer interface for mocking
type dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type mockDialer struct {
	called bool
}

func (m *mockDialer) DialAndSend(msg ...*gomail.Message) error {
	m.called = true
	return nil
}

func TestSendVerificationEmail(t *testing.T) {
	// Set environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "user")
	os.Setenv("SMTP_PASSWORD", "pass")

	// ToDO: insert mock here:
	//	originalNewDialer := newDialer
	//	defer func() { newDialer = originalNewDialer }()

	mock := &mockDialer{}
	// newDialer = func(host string, port int, username, password string) dialer {
	// 	return mock
	// }
	//
	SendVerificationEmail("test@example.com", "https://verify.com")

	if !mock.called {
		t.Error("Expected email to be sent")
	}
}
