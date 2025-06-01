package utils

import (
	"strings"
)

// MaskCardCode masks the card code for display/logging, showing only the first 2 and last 2 characters
func MaskCardCode(cardCode string) string {
	if len(cardCode) <= 4 {
		return "****"
	}
	
	if len(cardCode) <= 6 {
		// For short codes, show first and last character only
		return cardCode[:1] + strings.Repeat("*", len(cardCode)-2) + cardCode[len(cardCode)-1:]
	}
	
	// For longer codes, show first 2 and last 2 characters
	return cardCode[:2] + strings.Repeat("*", len(cardCode)-4) + cardCode[len(cardCode)-2:]
}

// MaskSensitiveString masks any sensitive string
func MaskSensitiveString(str string, showChars int) string {
	if len(str) <= showChars*2 {
		return strings.Repeat("*", len(str))
	}
	
	return str[:showChars] + strings.Repeat("*", len(str)-showChars*2) + str[len(str)-showChars:]
}

// SanitizeInput removes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim spaces
	input = strings.TrimSpace(input)
	
	return input
}

// IsValidCardCode validates card code format
func IsValidCardCode(cardCode string) bool {
	// Card code should be alphanumeric and 6-10 characters
	if len(cardCode) < 6 || len(cardCode) > 10 {
		return false
	}
	
	// Check if all characters are alphanumeric
	for _, ch := range cardCode {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')) {
			return false
		}
	}
	
	return true
}