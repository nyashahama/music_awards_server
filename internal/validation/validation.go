// Package validation
package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword follows OWASP and NIST 800-63B guidelines:
// - Minimum 8 characters (NIST recommends 8, max can be 64+)
// - Maximum 128 characters (practical limit)
// - No composition rules (no required character types)
// - Check against breached/common passwords
// - Allow all printable characters including spaces and unicode
func ValidatePassword(pw string) error {
	const (
		minLength = 8
		maxLength = 128
	)

	// Check length (use rune count for proper unicode handling)
	length := utf8.RuneCountInString(pw)

	if length < minLength {
		return errors.New("password must be at least 8 characters long")
	}

	if length > maxLength {
		return errors.New("password exceeds maximum length of 128 characters")
	}

	// Check for common/breached passwords
	if isCommonPassword(pw) {
		return errors.New("password is too common or has been found in data breaches - please choose a different password")
	}

	// Optional: Check for passwords that are just spaces
	if strings.TrimSpace(pw) == "" {
		return errors.New("password cannot be only whitespace")
	}

	return nil
}

// isCommonPassword checks against a list of commonly used/breached passwords
// In production, this should check against:
// - Have I Been Pwned API (https://haveibeenpwned.com/API/v3)
// - Local database of top 10k-100k most common passwords
// - Company/service-specific patterns (e.g., company name + year)
func isCommonPassword(pw string) bool {
	// Top commonly breached passwords (abbreviated list)
	// In production, use a comprehensive list (10k-100k entries) or API check
	commonPasswords := []string{
		"password", "12345678", "123456789", "1234567890",
		"qwerty", "abc123", "password1", "Password1",
		"11111111", "iloveyou", "welcome", "monkey",
		"dragon", "master", "sunshine", "princess",
		"letmein", "football", "shadow", "superman",
		"michael", "jennifer", "computer", "trustno1",
		"passw0rd", "admin", "user", "root",
	}

	// Case-insensitive check
	pwLower := strings.ToLower(pw)
	for _, common := range commonPasswords {
		if pwLower == strings.ToLower(common) {
			return true
		}
	}

	// Check for sequential patterns
	if isSequentialPattern(pw) {
		return true
	}

	return false
}

// isSequentialPattern detects simple sequential patterns like "abcdefgh" or "12345678"
func isSequentialPattern(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	// Check for repeated single character (e.g., "aaaaaaaa")
	if strings.Count(pw, string(pw[0])) == len(pw) {
		return true
	}

	// Check for simple numeric sequences
	sequences := []string{
		"01234567", "12345678", "23456789", "87654321",
		"abcdefgh", "qwertyui", "asdfghjk", "zxcvbnm",
	}

	pwLower := strings.ToLower(pw)
	for _, seq := range sequences {
		if strings.Contains(pwLower, seq) {
			return true
		}
	}

	return false
}
