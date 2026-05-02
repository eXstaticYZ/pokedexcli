package main

import "strings"

// cleanInput splits user input into words, lowercases them, and trims whitespace.
func cleanInput(text string) []string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return []string{}
	}

	// Convert to lowercase
	lowercase := strings.ToLower(trimmed)

	// Split by whitespace
	words := strings.Fields(lowercase)

	return words
}
