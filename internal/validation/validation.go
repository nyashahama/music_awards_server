package validation

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(pw string) error {
	n := len(pw)

	var hasLower, hasUpper, hasDigit bool
	repeatChanges := 0
	//repeatCount := 2 // track runs of repeating chars

	for i := 0; i < n; i++ {
		c := rune(pw[i])
		if unicode.IsLower(c) {
			hasLower = true
		}
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}

		// Count sequences of 3+ repeated characters
		if i >= 2 && pw[i] == pw[i-1] && pw[i-1] == pw[i-2] {
			repeatChanges++
		}
	}

	missingTypes := 0
	if !hasLower {
		missingTypes++
	}
	if !hasUpper {
		missingTypes++
	}
	if !hasDigit {
		missingTypes++
	}

	if n < 6 {
		// Need to insert enough to reach 6 characters and cover missing types
		steps := max(missingTypes, 6-n)
		return errors.New("Password is weak. Steps needed to strengthen: " + strconv.Itoa(steps))
	} else if n <= 20 {
		// Only replacement is needed
		steps := max(missingTypes, repeatChanges)
		if steps == 0 {
			return nil
		}
		return errors.New("Password is weak. Steps needed to strengthen: " + strconv.Itoa(steps))
	} else {
		// Need deletions
		over := n - 20
		// Use deletions to reduce repeating sequences
		steps := over + max(missingTypes, max(repeatChanges-over, 0))
		return errors.New("Password is weak. Steps needed to strengthen: " + strconv.Itoa(steps))
	}
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
