package main

import (
	"strings"
	"unicode"
)

// CleanWord converts the string to lowercase and removes punctuation and symbol characters.
func CleanWord(w string) string {
	w = strings.ToLower(w)
	var sb strings.Builder
	for _, r := range w {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

// IsValid returns true if the word has at least 3 characters.
// We use rune count to represent unicode character count correctly.
func IsValid(w string) bool {
	return len([]rune(w)) >= 3
}

// TokenizeLine splits a line into fields, cleans them, and filters valid ones.
func TokenizeLine(line string) []string {
	fields := strings.Fields(line)
	var words []string
	for _, f := range fields {
		cleaned := CleanWord(f)
		if IsValid(cleaned) {
			words = append(words, cleaned)
		}
	}
	return words
}
