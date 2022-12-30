package service

import (
	"testing"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "value 0",
			input:    0,
			expected: "0h0m",
		},
		{
			name:     "value length greater than 2",
			input:    25,
			expected: "25s",
		},
		{
			name:     "value length greater than 4",
			input:    451624,
			expected: "45h16m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Actual  %s", tt.expected, result)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "value greater than mb",
			input:    43534553,
			expected: "43.5MB",
		},
		{
			name:     "value greater than gb",
			input:    43534534234,
			expected: "43.5GB",
		},
		{
			name:     "value greater than tb",
			input:    41253453428734,
			expected: "41.3TB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSize(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Actual  %s", tt.expected, result)
			}
		})
	}
}
