package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsTimeBetween(t *testing.T) {
	ok := IsTimeBetween(time.Date(2023, time.December, 20, 16, 34, 10, 10, time.UTC), 1534, 2145)
	assert.Equal(t, ok, true)
	ok = IsTimeBetween(time.Date(2023, time.December, 20, 1, 34, 10, 10, time.UTC), 234, 2145)
	assert.Equal(t, ok, false)
}
func TestVerifyPassword(t *testing.T) {
	err := VerifyPassword("123456", "$2a$10$yjkyUPKwWo6IoSk1Nx85lujXWi0hpJ1omamBGmd4xtLFVvOVJXdh2")
	if err != nil {
		t.Error(err)
	}
}
