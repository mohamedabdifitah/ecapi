package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestProtectFields(t *testing.T) {
	exclude := ProtectFields("password", "devices")
	var Exclude = bson.D{
		{Key: "password", Value: 0},
		{Key: "devices", Value: 0},
	}
	assert.Equal(t, Exclude, exclude)
}
