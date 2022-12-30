package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	expected := "45h2m"
	result := formatDuration(450234)
	assert.Equal(t, expected, result)
}

func TestFormatSize(t *testing.T) {
	expected := "43.5MB"
	result := formatSize(43534553)
	assert.Equal(t, expected, result)
}
