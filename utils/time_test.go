package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTimeBetween(t *testing.T) {
	ok := IsTimeBetween(1534, 2145)
	assert.Equal(t, ok, false)
	ok = IsTimeBetween(234, 2145)
	assert.Equal(t, ok, false)
}
