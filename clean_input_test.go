package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", []string{}},
		{"single word", "hello", []string{"hello"}},
		{"multiple words", "hello world", []string{"hello", "world"}},
		{"mixed case", "Charmander Bulbasaur PIKACHU", []string{"charmander", "bulbasaur", "pikachu"}},
		{"leading/trailing spaces", "  hello world  ", []string{"hello", "world"}},
		{"multiple spaces", "hello   world", []string{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanInput(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %s, got %s", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
