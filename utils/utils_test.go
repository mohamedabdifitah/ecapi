package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsTimeBetween(t *testing.T) {
	ok := IsTimeBetween(time.Date(2023, time.December, 20, 16, 34, 10, 10, time.UTC), 15.34, 21.45)
	assert.Equal(t, ok, true)
	ok = IsTimeBetween(time.Date(2023, time.December, 20, 1, 34, 10, 10, time.UTC), 23.4, 21.45)
	assert.Equal(t, ok, false)
	ok = IsTimeBetween(time.Date(2023, time.December, 20, 0, 10, 0, 10, time.UTC), 0, 0)
	assert.Equal(t, ok, false)
	ok = IsTimeBetween(time.Date(2023, time.December, 20, 14, 00, 0, 10, time.UTC), 12.00, 23.00)
	assert.Equal(t, ok, true)
	ok = IsTimeBetween(time.Date(2023, time.December, 20, 23, 00, 0, 10, time.UTC), 1.00, 6.00)
	assert.Equal(t, ok, false)
}
func TestVerifyPassword(t *testing.T) {
	err := VerifyPassword("123456", "$2a$10$yjkyUPKwWo6IoSk1Nx85lujXWi0hpJ1omamBGmd4xtLFVvOVJXdh2")
	if err != nil {
		t.Error(err)
	}
}
func TestGenerateIDs(t *testing.T) {
	r := GenerateIDs(8)
	fmt.Println(r)
}
